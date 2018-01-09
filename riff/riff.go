package riff

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin/json"
	"io"
)

func (s *Server) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, "{")
	sortedNodes := s.Nodes.sort()
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
	digests = make(map[string]*Digest)
	for _, n := range s.Nodes {
		digest := &Digest{
			Id:       n.Id,
			SnapShot: n.SnapShot,
		}
		digests[digest.Id] = digest
	}
	d, _ := json.Marshal(digests)
	s.logger.Printf(infoRpcPrefix+"server %s send digests: %s\n", s.Self.Name, string(d))
	return
}

func (s *Server) MakeDiffNodes(digests Digests) (diff Nodes) {
	diff = make(map[string]*Node)
	for _, d := range digests {
		//find in server nodes
		n := s.Nodes[d.Id]
		if n == nil {
			//make an empty node for remote diff snap is empty
			empty := &Node{
				Id:      d.Id,
				Version: 0,
			}
			diff[d.Id] = empty
		} else {
			if d.SnapShot != n.SnapShot {
				diff[d.Id] = n
			}
		}
	}
	for _, n := range s.Nodes {
		if diff[n.Id] == nil && digests[n.Id] == nil {
			//add this server nodes
			diff[n.Id] = n
		}
	}
	d, _ := json.Marshal(diff)
	s.logger.Printf(infoRpcPrefix+"server %s send nodes: %s\n", s.Self.Name, string(d))
	return
}

func (s *Server) MergeDiff(diff Nodes) (reDiff Nodes) {
	reDiff = make(map[string]*Node)
	for _, d := range diff {
		n := s.Nodes[d.Id] //find in server nodes
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
				reDiff[n.Id] = n //shot out new version
			} else {
				if d.SnapShot == "" {
					//need update this
					reDiff[n.Id] = n
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
						reDiff[n.Id] = n
					}
				}
			}
		}
	}
	d, _ := json.Marshal(diff)
	r, _ := json.Marshal(reDiff)
	s.logger.Printf(infoRpcPrefix+"merge nodes: %s\n", string(d))
	s.logger.Printf(infoRpcPrefix+"return nodes: %s\n", string(r))
	s.Shutter()
	return
}

func (s *Server) Shutter() {
	h := sha1.New()
	io.WriteString(h, s.String())
	s.SnapShot = fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Server) AddNode(node *Node) *Node {
	if nd := s.Nodes[node.Id]; nd != nil {
		node = nd
	} else {
		s.Nodes[node.Id] = node
		node.Shutter()
	}
	return node
}

//func (s *Server) Link(node *Node, service *Service) {
//	s.AddNode(node).AddService(service)
//	s.Shutter()
//}
