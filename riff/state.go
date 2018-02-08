package riff

import (
	"fmt"
	"github.com/gimke/riff/api"
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
				node.State != api.StateAlive
		})
		if len(nodes) == 0 {
			addrs := strings.Split(s.config.Join, ",")
			for _, addr := range addrs {
				if addr != "" && addr != s.Self.Address() {
					if err := s.requestPeer(addr); err != nil {
						s.Logger.Printf(errorServerPrefix+"%v\n", err)
					} else {
						break
					}
				}
			}
		} else {
			for _, n := range nodes {
				if err := s.requestPeer(n.Address()); err != nil {
					s.SetStateWithShutter(n, api.StateSuspect)
					s.Logger.Printf(errorServerPrefix+"%v\n", err)
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
				node.State == api.StateAlive
		})

		for _, n := range nodes {
			if err := s.requestPeer(n.Address()); err != nil {
				if n.State == api.StateSuspect {
					s.SetStateWithShutter(n, api.StateDead)
					s.RemoveNodeDelay(n)
				}
				//s.logger.Printf(errorRpcPrefix+"%v\n", err)
			} else {
				s.SetStateWithShutter(n, api.StateAlive)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func (s *Server) fanoutLeave() {
	nodes := s.randomNodes(s.config.Fanout, func(node *Node) bool {
		return node.Name == s.Self.Name ||
			node.State != api.StateAlive
	})
	s.SetState(s.Self, api.StateSuspect)
	s.SetState(s.Self, api.StateDead)
	for _, n := range nodes {
		s.Logger.Printf(infoServerPrefix+"server %s send leave event to %s", s.Self.Name, n.Address())
		if err := s.requestLeave(n.Address()); err != nil {
			s.Logger.Printf(errorServerPrefix+"%v\n", err)
		}
	}
}

func (s *Server) requestLeave(peer string) error {
	client,err := api.NewClient(peer)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer client.Close()

	diff := []*Node{s.Self}

	var reDiff []*Node
	err = client.Call("Riff.PushDiff", diff, &reDiff)

	return err
}

func (s *Server) requestPeer(peer string) error {
	client,err := api.NewClient(peer)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer client.Close()

	var digests []*Digest
	err = client.Call("Riff.Request", s.SnapShot, &digests)
	if err != nil {
		return fmt.Errorf("peer: %s error: %v", peer, err)
	}
	//push diff
	if len(digests) != 0 {
		diff := s.MakeDiffNodes(digests)
		if len(diff) != 0 {
			var reDiff []*Node
			err = client.Call("Riff.PushDiff", diff, &reDiff)
			if len(reDiff) != 0 {
				s.MergeDiff(reDiff)
			}
		}
	}
	return nil
}
