package riff

import (
	"fmt"
	"github.com/gimke/riff/common"
	"net"
	"net/rpc"
	"time"
)

func (s *Server) stateFanout() {
	for {
		nodes := s.Nodes.randomNodes(s.config.Fanout, func(node *Node) bool {
			return node.Id == s.Self.Id ||
				node.State != stateAlive
		})
		if len(nodes) == 0 {
			//append join
			if s.config.Join != "" {
				nodes = append(nodes, s.config.Join)
			}
		}

		for _, peer := range nodes {
			if err := s.requestPeer(peer); err != nil {
				s.logger.Printf(errorRpcPrefix+"request peer error: %v", err)
			}
		}
		time.Sleep(3 * time.Second)
	}
}
func (s *Server) requestPeer(peer string) error {
	conn, err := net.DialTimeout("tcp", peer, time.Second*10)
	if err != nil {
		return fmt.Errorf("peer: %s error: %v", peer, err)
	}
	codec := common.NewGobClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	var digest []Node
	err = cmd.Call("Riff.Request", s.SnapShot, &digest)
	if err != nil {
		return fmt.Errorf("peer: %s error: %v", peer, err)
	}
	//push diff
	var diff []Node
	for _, n := range digest {
		remote := n
		my := s.Nodes[n.Id]
		if my == nil {
			empty := Node{
				Id:      n.Id,
				Version: 0,
			}
			diff = append(diff, empty)
		} else {
			if remote.SnapShot != my.SnapShot {
				diff = append(diff, *my)
			}
		}
	}
	var remoteDiff []Node
	err = cmd.Call("Riff.PushDiff", diff, &remoteDiff)
	s.MergeDiff(remoteDiff)
	s.Shutter()
	cmd.Close()
	return nil
}
