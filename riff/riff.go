package riff

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"time"
)

type Riff struct {
	Name       string
	DataCenter string
	Nodes
	Services
	SnapShort string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Create(config *Config) (*Riff, error) {
	riff := &Riff{
		Name:       config.Name,
		DataCenter: config.DataCenter,
		Nodes:      make(map[string]*Node),
		Services:   make(map[string]*Service),
	}
	return riff, nil
}

func (r *Riff) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, "{")
	sortedNodes := r.Nodes.sort()
	for i, nk := range sortedNodes {
		io.WriteString(buff, r.Nodes[nk].String())
		if i != len(sortedNodes)-1 {
			io.WriteString(buff, ",")
		}
	}
	io.WriteString(buff, "}")

	return buff.String()
}

func (r *Riff) Shutter() {
	h := sha1.New()
	io.WriteString(h, r.String())
	r.SnapShort = fmt.Sprintf("%x", h.Sum(nil))
}

func (r *Riff) AddNode(n *Node) *Node {
	if nd := r.Nodes[n.Name]; nd != nil {
		n = nd
	} else {
		r.Nodes[n.Name] = n
	}
	return n
}

func (r *Riff) Link(n *Node, s *Service) {
	r.AddNode(n).AddService(s)
}
