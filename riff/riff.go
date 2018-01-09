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
	sortedNodes := s.Nodes.Sort()
	for i, nk := range sortedNodes {
		//shutter the node
		//s.Nodes[nk].Shutter()
		io.WriteString(buff, s.Nodes[nk].SnapShot)
		if i != len(sortedNodes)-1 {
			io.WriteString(buff, ",")
		}
	}
	io.WriteString(buff, "}")

	return buff.String()
}

func (s *Server) MakeDigest() (digests Digests) {
	s.RLock()
	defer s.RUnlock()
	digests = make(map[string]*Digest)
	for _, n := range s.Nodes {
		digest := &Digest{
			Name:     n.Name,
			SnapShot: n.SnapShot,
		}
		digests[digest.Name] = digest
	}
	s.logger.Printf(infoRpcPrefix+"server %s send digests count: %d\n", s.Self.Name, len(digests))
	return
}

func (s *Server) MakeDiffNodes(digests Digests) (diff Nodes) {
	s.RLock()
	defer s.RUnlock()
	diff = make(map[string]*Node)
	for _, d := range digests {
		//find in server nodes
		n := s.Nodes[d.Name]
		if n == nil {
			//make an empty node for remote diff snap is empty
			empty := &Node{
				Name:    d.Name,
				Version: 0,
			}
			diff[d.Name] = empty
		} else {
			if d.SnapShot != n.SnapShot {
				diff[d.Name] = n
			}
		}
	}
	for _, n := range s.Nodes {
		if diff[n.Name] == nil && digests[n.Name] == nil {
			//add this server nodes
			diff[n.Name] = n
		}
	}
	s.logger.Printf(infoRpcPrefix+"server %s send nodes count: %d\n", s.Self.Name, len(diff))
	return
}

func (s *Server) MergeDiff(diff Nodes) (reDiff Nodes) {
	reDiff = make(map[string]*Node)
	for _, d := range diff {
		n := s.Nodes[d.Name] //find in server nodes
		if n == nil {
			d.IsSelf = false //remove is self
			s.AddNode(d)     //if not find then add node
		} else {
			if d.IsSelf {
				//if remote node is self then overwrite server node
				v := n.Version + 1
				*n = *d
				n.IsSelf = false
				n.Version = v
				n.Shutter()
				reDiff[n.Name] = n //shot out new version
			} else {
				if d.SnapShot == "" {
					//need update this
					reDiff[n.Name] = n
				} else {
					if d.Version > n.Version {
						if n.IsSelf {
							//only update version
							n.Version = d.Version
							n.Shutter()
						} else {
							*n = *d
						}
					} else if d.Version != n.Version {
						//take my
						reDiff[n.Name] = n
					}
				}
			}
		}
	}
	s.Shutter()
	return
}

func (s *Server) Shutter() {
	s.Lock()
	defer s.Unlock()
	h := sha1.New()
	io.WriteString(h, s.String())
	s.SnapShot = fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Server) AddNode(node *Node) *Node {
	s.Lock()
	defer s.Unlock()
	if nd := s.Nodes[node.Name]; nd != nil {
		node = nd
	} else {
		s.Nodes[node.Name] = node
		node.Shutter()
	}
	return node
}

//func (s *Server) Link(node *Node, service *Service) {
//	s.AddNode(node).AddService(service)
//	s.Shutter()
//}
