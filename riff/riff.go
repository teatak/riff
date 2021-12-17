package riff

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/teatak/riff/api"
)

func (s *Server) Shutter() {
	h := sha1.New()
	io.WriteString(h, s.String())
	s.SnapShot = fmt.Sprintf("%x", h.Sum(nil))
	s.Logger.Printf(infoServerPrefix+"snapshot now is: %s\n", s.SnapShot)
}

func (s *Server) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, "{")
	keys := s.Keys()
	for i, nk := range keys {
		n := s.GetNode(nk)
		if n != nil && n.State != api.StateDead {
			io.WriteString(buff, n.SnapShot)
			if i != len(keys)-1 {
				io.WriteString(buff, ",")
			}
		}
	}
	io.WriteString(buff, "}")
	return buff.String()
}

func (s *Server) MakeDigest() (digests []*Digest) {
	digests = make([]*Digest, 0)
	s.Range(func(key string, n *Node) bool {
		digest := &Digest{
			Name:     n.Name,
			SnapShot: n.SnapShot,
		}
		digests = append(digests, digest)
		return true
	})
	return
}

func (s *Server) MakeDiffNodes(digests []*Digest) (diff []*Node) {
	diff = make([]*Node, 0)
	keysDiff := make(map[string]bool)
	keysDigest := make(map[string]bool)
	for _, d := range digests {
		keysDigest[d.Name] = true
		//find in server nodes
		n := s.GetNode(d.Name)
		if n == nil {
			//make an empty node for remote diff snap is empty
			empty := &Node{
				Name:    d.Name,
				Version: 0,
			}
			diff = append(diff, empty)
			keysDiff[d.Name] = true
		} else {
			if d.SnapShot != n.SnapShot {
				diff = append(diff, n)
				keysDiff[n.Name] = true
			}
		}
	}

	s.Range(func(key string, node *Node) bool {
		if !keysDiff[node.Name] && !keysDigest[node.Name] {
			diff = append(diff, node)
		}
		return true
	})
	return
}

func (s *Server) MergeDiff(diffs []*Node) (reDiffs []*Node) {
	reDiffs = make([]*Node, 0)
	count := 0
	for _, d := range diffs {
		n := s.GetNode(d.Name) //find in server nodes
		if n == nil {          //if not found in server then add node
			if d.State != api.StateDead {
				//exclude dead node
				d.IsSelf = false //remove is self
				s.AddNode(d)     //if not find then add node
				count++
			}
		} else {

			// if found in node map
			if d.SnapShot == n.SnapShot {
				// continue if have same snap
				continue
			}
			if d.SnapShot == "" {
				//other cluster will send empty snap node for query this node
				reDiffs = append(reDiffs, n)
				continue
			}

			var merged bool
			var reDiff *Node

			switch d.IsSelf {
			case true:
				merged, reDiff = s.trueNode(d, n)
			case false:
				merged, reDiff = s.gossipNode(d, n)
			}
			if reDiff != nil {
				reDiffs = append(reDiffs, reDiff)
			}
			if merged {
				count++
			}
		}
	}
	s.Logger.Printf(infoServerPrefix+"server %s merge %d nodes return %d nodes\n", s.Self.Name, count, len(reDiffs))
	s.Shutter()
	return
}

// it's real true node state
func (s *Server) trueNode(d, n *Node) (merged bool, reDiff *Node) {
	//if remote node is self then overwrite server node
	switch d.State {
	case api.StateAlive:
		if d.VersionGet() == 0 { //if d is new online
			//v := s.SetState(n, api.StateAlive)
			v := n.VersionGet()
			s.ExchangeNode(n, d)
			n.VersionSet(v)
			v = s.SetState(n, api.StateAlive)
			if v > 1 {
				reDiff = n //shot out my version
			}
		} else {
			//if remote node service changes .... take remote node
			s.ExchangeNode(n, d)
		}
		merged = true
	case api.StateDead:
		if n.State != api.StateDead {
			s.ExchangeNode(n, d)
			s.RemoveNodeDelay(n)
			merged = true
		}
	}
	return
}

func (s *Server) gossipNode(d, n *Node) (merged bool, reDiff *Node) {
	if d.VersionGet() > n.VersionGet() {
		if n.IsSelf {
			//only update version
			n.VersionSet(d.Version)
			n.Shutter()
		} else {
			if n.State != api.StateDead && d.State == api.StateDead {
				s.ExchangeNode(n, d)
				s.RemoveNodeDelay(n)
			} else {
				s.ExchangeNode(n, d)
			}
		}
		merged = true
	} else if d.VersionGet() != n.VersionGet() {
		//take my node
		reDiff = n
	}
	return
}

func (s *Server) ExchangeNode(n, d *Node) {
	//discover diff
	if n.SnapShot == d.SnapShot || n.Name != d.Name {
		return
	}

	s.walkService(n, d)

	if n.IsSelf {
		n = d
	} else {
		n = d
		n.IsSelf = false
	}

	s.watch.Dispatch(WatchParam{
		Name:      n.Name,
		WatchType: NodeChanged,
	})
}

func (s *Server) walkService(n, d *Node) {
	walked := make(map[string]bool)
	for _, sv := range n.Services {
		//find in diff service
		walked[sv.Name] = true
		if d.Services[sv.Name] == nil {
			s.watch.Dispatch(WatchParam{
				Name:      sv.Name,
				WatchType: ServiceChanged,
			})
		} else {
			if d.Services[sv.Name].State != sv.State {
				s.watch.Dispatch(WatchParam{
					Name:      sv.Name,
					WatchType: ServiceChanged,
				})
			}
			if d.Services[sv.Name].Progress.Current != sv.Progress.Current {
				s.watch.Dispatch(WatchParam{
					Name:      sv.Name,
					WatchType: ServiceChanged,
				})
			}
		}
	}
	for _, sv := range d.Services {
		if !walked[sv.Name] {
			s.watch.Dispatch(WatchParam{
				Name:      sv.Name,
				WatchType: ServiceChanged,
			})
		}
	}
}
