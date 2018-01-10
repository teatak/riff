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
		node, _ := s.Load(nk)
		if node != nil {
			io.WriteString(buff, node.(*Node).SnapShot)
			if i != len(sortedNodes)-1 {
				io.WriteString(buff, ",")
			}
		}
	}
	io.WriteString(buff, "}")

	return buff.String()
}

func (s *Server) MakeDigest() (digests []*Digest) {
	digests = make([]*Digest, 0)
	s.Range(func(key, value interface{}) bool {
		n := value.(*Node)
		digest := &Digest{
			Name:     n.Name,
			SnapShot: n.SnapShot,
		}
		digests = append(digests, digest)
		return true
	})
	//b,_ := json.Marshal(digests)
	//fmt.Println(string(b))
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
		node, _ := s.Load(d.Name)
		if node == nil {
			//make an empty node for remote diff snap is empty
			empty := &Node{
				Name:    d.Name,
				Version: 0,
			}
			diff = append(diff, empty)
			keysDiff[d.Name] = true
		} else {
			n := node.(*Node)
			if d.SnapShot != n.SnapShot {
				diff = append(diff, n)
				keysDiff[n.Name] = true
			}
		}
	}

	s.Range(func(key, value interface{}) bool {
		n := value.(*Node)
		if !keysDiff[n.Name] && !keysDigest[n.Name] {
			diff = append(diff, n)
		}
		return true
	})

	s.logger.Printf(infoRpcPrefix+"server %s get %d digests send %d nodes\n", s.Self.Name, len(digests), len(diff))
	return
}

func (s *Server) MergeDiff(diff []*Node) (reDiff []*Node) {
	reDiff = make([]*Node, 0)
	for _, d := range diff {
		node, _ := s.Load(d.Name) //find in server nodes
		if node == nil {
			d.IsSelf = false //remove is self
			s.AddNode(d)     //if not find then add node
		} else {
			n := node.(*Node)
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
				n.Shutter()
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

func (s *Server) AddNode(node *Node) {
	//if nd := s.Nodes[node.Name]; nd != nil {
	//	node = nd
	//} else {
	//	s.Nodes[node.Name] = node
	//	node.Shutter()
	//}
	node.Shutter()
	s.Store(node.Name, node)
}

//func (s *Server) Link(node *Node, service *Service) {
//	s.AddNode(node).AddService(service)
//	s.Shutter()
//}
