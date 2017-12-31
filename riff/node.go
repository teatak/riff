package riff

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strconv"
)

type Nodes map[string]*Node

type Node struct {
	Services
	Name       string
	DataCenter string
	SnapShort  string
}

func (ns *Nodes) sort() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *ns {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (n *Node) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, n.Name+":{")
	//write service name and version
	sortedServices := n.Services.sort()
	for i, sk := range sortedServices {
		s := n.Services[sk]
		io.WriteString(buff, s.Name+":{"+s.Address+","+strconv.FormatUint(s.Version, 10)+"}")
		if i != len(sortedServices)-1 {
			io.WriteString(buff, ",")
		}
	}
	io.WriteString(buff, "}")
	return buff.String()
}

func (n *Node) Shutter() {
	h := sha1.New()
	io.WriteString(h, n.String())
	n.SnapShort = fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(n.Name, "SnapShort", n.SnapShort)
}

func (n *Node) AddService(s *Service) {
	if n.Services == nil {
		n.Services = make(map[string]*Service)
	}
	n.Services[s.Name] = s
	n.Shutter()
}
