package riff

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}


func (s *Server) toString() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, "{")
	sortedNodes := s.Nodes.sort()
	for i, nk := range sortedNodes {
		//shutter the node
		s.Nodes[nk].Shutter()
		io.WriteString(buff, s.Nodes[nk].toString())
		if i != len(sortedNodes)-1 {
			io.WriteString(buff, ",")
		}
	}
	io.WriteString(buff, "}")

	return buff.String()
}

func (s *Server) Shutter() {
	h := sha1.New()
	io.WriteString(h, s.toString())
	s.SnapShort = fmt.Sprintf("%x", h.Sum(nil))
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
