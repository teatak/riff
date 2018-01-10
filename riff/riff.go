package riff

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
)

func (s *Server) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, "{")
	keys := s.Nodes.Keys()
	for i, nk := range keys {
		n := s.Get(nk)
		if n != nil && n.State != stateDead {
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
	s.logger.Printf(infoRpcPrefix+"server %s send %d digests\n", s.Self.Name, len(digests))
	return
}

func (s *Server) MakeDiffNodes(digests []*Digest) (diff []*Node) {
	diff = make([]*Node, 0)
	keysDiff := make(map[string]bool)
	keysDigest := make(map[string]bool)
	for _, d := range digests {
		keysDigest[d.Name] = true
		//find in server nodes
		n := s.Get(d.Name)
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

	s.Range(func(key string,node *Node) bool {
		if !keysDiff[node.Name] && !keysDigest[node.Name] {
			diff = append(diff, node)
		}
		return true
	})

	s.logger.Printf(infoRpcPrefix+"server %s get %d digests send %d nodes\n", s.Self.Name, len(digests), len(diff))
	return
}

func (s *Server) MergeDiff(diff []*Node) (reDiff []*Node) {
	reDiff = make([]*Node, 0)
	for _, d := range diff {
		n := s.Get(d.Name) //find in server nodes
		if n == nil {
			d.IsSelf = false //remove is self
			s.SetNode(d)     //if not find then add node
		} else {
			if d.SnapShot == n.SnapShot {
				continue
			}
			if d.SnapShot == "" {
				//need update this
				reDiff = append(reDiff, n)
				continue
			}
			if d.IsSelf {
				//if remote node is self then overwrite server node
				v := n.Version + 1
				*n = *d
				n.IsSelf = false
				n.Version = v
				s.SetNode(n)
				reDiff = append(reDiff, n)
				//reDiff[n.Name] = n //shot out new version
			} else {
				if d.Version > n.Version {
					if n.IsSelf {
						//only update version
						n.Version = d.Version
						n.Shutter()
					} else {
						*n = *d
						s.SetNode(n)
					}
				} else if d.Version != n.Version {
					//take my node
					reDiff = append(reDiff, n)
				}
			}
		}
	}
	s.logger.Printf(infoRpcPrefix+"server %s merge %d nodes return %d nodes\n", s.Self.Name, len(diff), len(reDiff))
	s.Shutter()
	return
}

func (s *Server) Shutter() {
	h := sha1.New()
	io.WriteString(h, s.String())
	s.SnapShot = fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Server) SetNode(node *Node) {
	node.Shutter()
	s.Set(node)
	if node.State == stateDead {
		node.Dead(s)
	}
}