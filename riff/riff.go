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
		//s.Nodes[nk].Shutter()
		io.WriteString(buff, s.Nodes[nk].SnapShot)
		if i != len(sortedNodes)-1 {
			io.WriteString(buff, ",")
		}
	}
	io.WriteString(buff, "}")

	return buff.String()
}

func (s *Server) MakeDigest() (digest []Node) {
	digest = make([]Node, 0)
	for _, n := range s.Nodes {
		new := *n
		new.Services = nil
		digest = append(digest, new)
	}
	return
}

func (s *Server) MergeDiff(diff []Node) []Node {
	nodes := make([]Node, 0)
	s.logger.Printf(infoRpcPrefix+"merge nodes len: %d", len(diff))
	for _, remote := range diff {
		myNode := s.Nodes[remote.Id]
		if myNode != nil {
			//find node
			if remote.IsSelf {
				//overwrite if remote node is self node
				//over write version
				version := myNode.Version
				*myNode = remote
				myNode.IsSelf = false
				myNode.Version = version
				myNode.Shutter()
				nodes = append(nodes, *myNode) //shot out new version
			} else {
				if remote.Version > myNode.Version {
					if myNode.IsSelf {
						//only update version
						myNode.Version = remote.Version
					} else {
						*myNode = remote
					}
					myNode.Shutter()
				} else {
					nodes = append(nodes, *myNode)
				}
			}
		} else {
			//not find node
			s.AddNode(&remote)
			s.Shutter()
		}
	}
	return nodes
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
