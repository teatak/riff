package riff

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Nodes map[string]*Node

type stateType int

const (
	stateAlive stateType = iota
	stateSuspect
	stateDead
)

func (s stateType) String() string {
	switch s {
	case stateAlive:
		return "Alive"
		break
	case stateSuspect:
		return "Suspect"
		break
	case stateDead:
		return "Dead"
		break
	}
	return "Unknow"
}

type Node struct {
	Services
	Name        string
	DataCenter  string
	IP          net.IP
	Port        int
	State       stateType // Current state
	StateChange time.Time // Time last state change happened
	SnapShort   string

	nodeLock sync.RWMutex
}

func (n *Node) Address() string {
	return net.JoinHostPort(n.IP.String(), strconv.Itoa(int(n.Port)))
}

func (ns *Nodes) sort() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *ns {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (n *Node) toString() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, n.Name+":{")
	//write service name and version
	sortedServices := n.Services.sort()
	for i, sk := range sortedServices {
		s := n.Services[sk]
		io.WriteString(buff, s.Name+":{"+s.Address()+","+strconv.FormatUint(s.Version, 10)+"}")
		if i != len(sortedServices)-1 {
			io.WriteString(buff, ",")
		}
	}
	io.WriteString(buff, "}")
	return buff.String()
}

func (n *Node) Shutter() {
	h := sha1.New()
	io.WriteString(h, n.toString())
	n.SnapShort = fmt.Sprintf("%x", h.Sum(nil))
}

func (n *Node) AddService(s *Service) {
	if n.Services == nil {
		n.Services = make(map[string]*Service)
	}
	n.Services[s.Name] = s
	n.Shutter()
}