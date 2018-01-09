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
			return node.Name == s.Self.Name ||
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
				s.logger.Printf(errorRpcPrefix+"request peer error: %v\n", err)
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}
func (s *Server) requestPeer(peer string) error {
	conn, err := net.DialTimeout("tcp", peer, time.Second*10)
	if err != nil {
		return fmt.Errorf("peer: %s error: %v", peer, err)
	}
	codec := common.NewGobClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	defer cmd.Close()

	var digests Digests
	err = cmd.Call("Riff.Request", s.SnapShot, &digests)
	if err != nil {
		return fmt.Errorf("peer: %s error: %v", peer, err)
	}
	//push diff
	if len(digests) != 0 {
		diff := s.MakeDiffNodes(digests)
		if len(diff) != 0 {
			var reDiff Nodes
			err = cmd.Call("Riff.PushDiff", diff, &reDiff)
			if len(reDiff) != 0 {
				s.MergeDiff(reDiff)
			}
		}
	}
	return nil
}
