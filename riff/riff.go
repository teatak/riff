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
	sortedNodes := s.Nodes.sort()
	for i, nk := range sortedNodes {
		//shutter the node
		s.Nodes[nk].Shutter()
		io.WriteString(buff, s.Nodes[nk].String())
		if i != len(sortedNodes)-1 {
			io.WriteString(buff, ",")
		}
	}
	io.WriteString(buff, "}")

	return buff.String()
}

func (s *Server) MakeDigest() (digest Nodes){
	digest = make(map[string]*Node)
	for _,n := range s.Nodes {
		digest[n.Name] = n
		//remove services
		digest[n.Name].Services = nil
	}
	return
}

func (s *Server) Shutter() {
	h := sha1.New()
	io.WriteString(h, s.String())
	s.SnapShot = fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Server) AddNode(node *Node) *Node {
	if nd := s.Nodes[node.Name]; nd != nil {
		node = nd
	} else {
		s.Nodes[node.Name] = node
	}
	s.Shutter()
	return node
}

func (s *Server) Link(node *Node, service *Service) {
	s.Shutter()
	s.AddNode(node).AddService(service)
}
