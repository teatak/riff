package riff

import (
	"fmt"
	"github.com/gimke/riff/api"
	"net"
	"net/rpc"
	"strings"
	"time"
)

func (s *Server) fanoutNodes() {
	for {
		select {
		case <-s.ShutdownCh:
			return
		default:
		}
		nodes := s.randomNodes(s.config.Fanout, func(node *Node) bool {
			return node.Name == s.Self.Name ||
				node.State != stateAlive
		})
		if len(nodes) == 0 {
			addrs := strings.Split(s.config.Join, ",")
			for _, addr := range addrs {
				if addr != "" && addr != s.Self.Address() {
					if err := s.requestPeer(addr); err != nil {
						s.Logger.Printf(errorRpcPrefix+"%v\n", err)
					} else {
						break
					}
				}
			}
		} else {
			for _, n := range nodes {
				if err := s.requestPeer(n.Address()); err != nil {
					s.SetStateWithShutter(n, stateSuspect)
					s.Logger.Printf(errorRpcPrefix+"%v\n", err)
				}
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (s *Server) fanoutDeadNodes() {
	for {
		select {
		case <-s.ShutdownCh:
			return
		default:
		}
		nodes := s.randomNodes(1, func(node *Node) bool {
			return node.Name == s.Self.Name ||
				node.State == stateAlive
		})

		for _, n := range nodes {
			if err := s.requestPeer(n.Address()); err != nil {
				if n.State == stateSuspect {
					s.SetStateWithShutter(n, stateDead)
					s.RemoveNodeDelay(n)
				}
				//s.logger.Printf(errorRpcPrefix+"%v\n", err)
			} else {
				s.SetStateWithShutter(n, stateAlive)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func (s *Server) fanoutLeave() {
	nodes := s.randomNodes(s.config.Fanout, func(node *Node) bool {
		return node.Name == s.Self.Name ||
			node.State != stateAlive
	})
	s.SetState(s.Self, stateSuspect)
	s.SetState(s.Self, stateDead)
	for _, n := range nodes {
		s.Logger.Printf(infoRpcPrefix+"server %s send leave event to %s", s.Self.Name, n.Address())
		if err := s.requestLeave(n.Address()); err != nil {
			s.Logger.Printf(errorRpcPrefix+"%v\n", err)
		}
	}
}

func (s *Server) requestLeave(peer string) error {
	conn, err := net.DialTimeout("tcp", peer, time.Second*10)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	codec := api.NewGobClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	defer cmd.Close()

	diff := []*Node{s.Self}

	var reDiff []*Node
	err = cmd.Call("Riff.PushDiff", diff, &reDiff)

	return err
}

func (s *Server) requestPeer(peer string) error {
	conn, err := net.DialTimeout("tcp", peer, time.Second*10)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	codec := api.NewGobClientCodec(conn)
	cmd := rpc.NewClientWithCodec(codec)
	defer cmd.Close()

	var digests []*Digest
	err = cmd.Call("Riff.Request", s.SnapShot, &digests)
	if err != nil {
		return fmt.Errorf("peer: %s error: %v", peer, err)
	}
	//push diff
	if len(digests) != 0 {
		diff := s.MakeDiffNodes(digests)
		if len(diff) != 0 {
			var reDiff []*Node
			err = cmd.Call("Riff.PushDiff", diff, &reDiff)
			if len(reDiff) != 0 {
				s.MergeDiff(reDiff)
			}
		}
	}
	return nil
}
