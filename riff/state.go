package riff

import (
	"fmt"
	"github.com/gimke/riff/common"
	"net"
	"net/rpc"
	"strings"
	"time"
)

func (s *Server) fanoutNodes() {
	for {
		nodes := s.Nodes.randomNodes(s.config.Fanout, func(node *Node) bool {
			return node.Name == s.Self.Name ||
				node.State != stateAlive
		})
		if len(nodes) == 0 {
			addrs := strings.Split(s.config.Join, ",")
			for _, addr := range addrs {
				if addr != "" && addr != s.Self.Address() {
					if err := s.requestPeer(addr); err != nil {
						s.logger.Printf(errorRpcPrefix+"request peer error: %v\n", err)
					}
				}
			}
		} else {
			for _, peer := range nodes {
				if err := s.requestPeer(peer.Address()); err != nil {
					peer.State = stateSuspect
					peer.Version++
					peer.Shutter()
					s.Shutter()
					s.logger.Printf(errorRpcPrefix+"request peer error: %v\n", err)
				}
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (s *Server) fanoutDeadNodes() {
	for {
		nodes := s.Nodes.randomNodes(1, func(node *Node) bool {
			return node.Name == s.Self.Name ||
				node.State == stateAlive
		})

		for _, peer := range nodes {
			if err := s.requestPeer(peer.Address()); err != nil {
				if peer.State == stateSuspect {
					peer.State = stateDead
					peer.Version++
					peer.Shutter()
					s.Shutter()
				}
				s.logger.Printf(errorRpcPrefix+"request dead peer error: %v\n", err)
			} else {
				peer.State = stateAlive
				peer.Version++
				peer.Shutter()
				s.Shutter()
			}
		}

		time.Sleep(2 * time.Second)
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
