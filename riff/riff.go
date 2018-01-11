package riff

import (
	"bytes"
	"io"
)

func (s *Server) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, "{")
	keys := s.Keys()
	for i, nk := range keys {
		n := s.GetNode(nk)
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
	s.Logger.Printf(infoNodePrefix+"server %s send %d digests\n", s.Self.Name, len(digests))
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

	s.Logger.Printf(infoNodePrefix+"server %s get %d digests send %d nodes\n", s.Self.Name, len(digests), len(diff))
	return
}

func (s *Server) MergeDiff(diff []*Node) (reDiff []*Node) {
	reDiff = make([]*Node, 0)
	for _, d := range diff {
		n := s.GetNode(d.Name) //find in server nodes
		if n == nil {
			if d.State != stateDead {
				//exclude dead node
				d.IsSelf = false //remove is self
				s.SetNode(d)     //if not find then add node
			}
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
				v := n.VersionInc()
				//n.Version + 1
				*n = *d
				n.IsSelf = false
				n.Witness(v)
				s.SetNode(n)
				reDiff = append(reDiff, n)
				//reDiff[n.Name] = n //shot out new version
			} else {
				if d.VersionGet() > n.VersionGet() {
					if n.IsSelf {
						//only update version
						n.Witness(d.Version)
						n.Shutter()
					} else {
						*n = *d
						s.SetNode(n)
					}
				} else if d.VersionGet() != n.VersionGet() {
					//take my node
					reDiff = append(reDiff, n)
				}
			}
		}
	}
	s.Logger.Printf(infoNodePrefix+"server %s merge %d nodes return %d nodes\n", s.Self.Name, len(diff), len(reDiff))
	s.Shutter()
	return
}
