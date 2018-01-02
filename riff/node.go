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

func (ns *Nodes) sort() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *ns {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

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
	Id          string
	Name        string
	DataCenter  string
	IP          string
	Port        int
	State       stateType // Current state
	StateChange time.Time // Time last state change happened
	SnapShot   string

	nodeLock sync.RWMutex
}

func (n *Node) Address() string {
	return net.JoinHostPort(n.IP, strconv.Itoa(int(n.Port)))
}

func (n *Node) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, n.Id+strconv.Itoa(n.StateChange.Nanosecond())+":{")
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
	io.WriteString(h, n.String())
	n.SnapShot = fmt.Sprintf("%x", h.Sum(nil))
}

func (n *Node) AddService(s *Service) {
	if n.Services == nil {
		n.Services = make(map[string]*Service)
	}
	n.Services[s.Name] = s
	n.Shutter()
}
