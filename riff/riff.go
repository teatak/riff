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
			if d.State != stateDead {
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
				break
			case false:
				merged, reDiff = s.gossipNode(d, n)
				break
			}
			if reDiff != nil {
				reDiffs = append(reDiffs, reDiff)
			}
			if merged {
				count++
			}
		}
	}
	s.Logger.Printf(infoNodePrefix+"server %s merge %d nodes return %d nodes\n", s.Self.Name, count, len(reDiffs))
	s.Shutter()
	return
	//reDiff = make([]*Node, 0)
	//for _, d := range diff {
	//	n := s.GetNode(d.Name) //find in server nodes
	//	if n == nil {
	//		if d.State != stateDead {
	//			//exclude dead node
	//			d.IsSelf = false //remove is self
	//			s.SetNode(d)     //if not find then add node
	//		}
	//	} else {
	//		if d.SnapShot == n.SnapShot {
	//			continue
	//		}
	//		if d.SnapShot == "" {
	//			//need update this
	//			reDiff = append(reDiff, n)
	//			continue
	//		}
	//		if d.IsSelf {
	//			//if remote node is self then overwrite server node
	//			switch d.State {
	//			case stateAlive:
	//				s.SetState(n, stateAlive)
	//				break
	//			case stateDead:
	//				n.State = stateDead
	//				n.VersionSet(d.Version)
	//				s.SetNode(n)
	//				s.RemoveNode(n)
	//				break
	//			}
	//			reDiff = append(reDiff, n)
	//			//reDiff[n.Name] = n //shot out new version
	//		} else {
	//			if d.VersionGet() > n.VersionGet() {
	//				if n.IsSelf {
	//					//only update version
	//					n.VersionSet(d.Version)
	//					n.Shutter()
	//				} else {
	//					*n = *d
	//					s.SetNode(n)
	//				}
	//			} else if d.VersionGet() != n.VersionGet() {
	//				//take my node
	//				reDiff = append(reDiff, n)
	//			}
	//		}
	//	}
	//}
	//s.Logger.Printf(infoNodePrefix+"server %s merge %d nodes return %d nodes\n", s.Self.Name, len(diff), len(reDiff))
	//s.Shutter()
	//return
}

// it's real true node state
func (s *Server) trueNode(d, n *Node) (merged bool, reDiff *Node) {
	//if remote node is self then overwrite server node
	switch d.State {
	case stateAlive:
		if d.VersionGet() == 0 { //if d is new online
			v := s.SetState(n, stateAlive)
			if v > 1 {
				reDiff = n //shot out my version
			}
		} else {
			//if remote node service changes .... take remote node
			*n = *d
		}
		merged = true
		break
	case stateDead:
		if n.State != stateDead {
			*n = *d
			s.RemoveNodeDelay(n)
			merged = true
		}
		break
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
			if n.State != stateDead && d.State == stateDead {
				*n = *d
				s.RemoveNodeDelay(n)
			} else {
				*n = *d
			}
		}
		merged = true
	} else if d.VersionGet() != n.VersionGet() {
		//take my node
		reDiff = n
	}
	return
}
