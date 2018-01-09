package riff

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/gimke/riff/common"
	"io"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Nodes map[string]*Node
type Digests map[string]*Digest
type Digest struct {
	Name     string
	SnapShot string
}

func (ns *Nodes) sort() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *ns {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (ns *Nodes) randomNodes(fanout int, filterFn func(*Node) bool) []string {
	nodes := ns.sort()
	n := len(nodes)
	kNodes := make([]*Node, 0, fanout)
	KString := make([]string, 0, fanout)
OUTER:
	for i := 0; i < 3*n && len(kNodes) < fanout; i++ {
		idx := common.RandomNumber(n)
		node := (*ns)[nodes[idx]]
		//filter nodes
		if filterFn != nil && filterFn(node) {
			continue OUTER
		}
		// Check if we have this node already
		for j := 0; j < len(kNodes); j++ {
			if node == kNodes[j] {
				continue OUTER
			}
		}
		kNodes = append(kNodes, node)
	}

	for i := 0; i < len(kNodes); i++ {
		KString = append(KString, kNodes[i].Address())
	}
	return KString
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
	Name        string
	DataCenter  string
	IP          string
	Port        int
	Version     uint64
	State       stateType // Current state
	StateChange time.Time // Time last state change happened
	SnapShot    string
	IsSelf      bool
	nodeLock    sync.RWMutex
}

func (n *Node) Address() string {
	return net.JoinHostPort(n.IP, strconv.Itoa(int(n.Port)))
}

func (n *Node) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, n.Name+":"+strconv.Itoa(n.StateChange.Nanosecond())+":{")
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
