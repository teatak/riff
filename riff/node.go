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

type Nodes struct {
	nodes sync.Map
}

type Digest struct {
	Name     string
	SnapShot string
}

func (ns *Nodes) Keys() []string {
	var keys = make([]string, 0, 0)
	ns.nodes.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	sort.Strings(keys)
	return keys
}

func (ns *Nodes) Slice() []*Node {
	keys := ns.Keys()
	nodes := make([]*Node, 0)
	for _, key := range keys {
		if n, _ := ns.nodes.Load(key); n != nil {
			nodes = append(nodes, n.(*Node))
		}
	}
	return nodes
}
func (ns *Nodes) Range(f func(string, *Node) bool) {
	ns.nodes.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Node))
	})
}
func (ns *Nodes) Get(key string) *Node {
	if n, _ := ns.nodes.Load(key); n != nil {
		return n.(*Node)
	}
	return nil
}

func (ns *Nodes) Set(value *Node) {
	ns.nodes.Store(value.Name, value)
}

func (ns *Nodes) randomNodes(fanout int, filterFn func(*Node) bool) []*Node {
	nodes := ns.Keys()
	n := len(nodes)
	RNodes := make([]*Node, 0, fanout)
OUTER:
	for i := 0; i < 3*n && len(RNodes) < fanout; i++ {
		idx := common.RandomNumber(n)
		n := ns.Get(nodes[idx])
		if n == nil {
			continue OUTER
		} //filter nodes
		if filterFn != nil && filterFn(n) {
			continue OUTER
		}
		// Check if we have this node already
		for j := 0; j < len(RNodes); j++ {
			if n == RNodes[j] {
				continue OUTER
			}
		}
		RNodes = append(RNodes, n)
	}

	return RNodes
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
	Name        string
	DataCenter  string
	IP          string
	Port        int
	Version     uint64
	State       stateType // Current state
	StateChange time.Time // Time last state change happened
	SnapShot    string
	Services
	IsSelf    bool
	nodeLock  sync.RWMutex
	timer     *time.Timer
	timeoutFn func()
}

func (n *Node) Address() string {
	return net.JoinHostPort(n.IP, strconv.Itoa(int(n.Port)))
}

func (n *Node) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, n.Name+":"+strconv.FormatUint(n.Version, 10)+":{")
	//write service name and version
	keys := n.Services.Keys()
	for i, sk := range keys {
		s := n.Services[sk]
		io.WriteString(buff, s.Name+":{"+s.Address()+","+strconv.FormatUint(s.Version, 10)+"}")
		if i != len(keys)-1 {
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

func (n *Node) Suspect() {
	n.State = stateSuspect
	n.Version++
	n.Shutter()
}

func (n *Node) Dead(s *Server) {
	if n.timeoutFn == nil {
		n.State = stateDead
		n.Version++
		n.Shutter()
		n.timeoutFn = func() {
			//delete this node
			if n.State == stateDead {
				s.Nodes.nodes.Delete(n.Name)
				s.Shutter()
			}
		}
		timeout := 10*time.Second
		n.timer = time.AfterFunc(timeout, n.timeoutFn)
	}
}

func (n *Node) Alive() {
	n.State = stateAlive
	n.Shutter()
}

func (n *Node) AddService(s *Service) {
	if n.Services == nil {
		n.Services = make(map[string]*Service)
	}
	n.Services[s.Name] = s
	n.Shutter()
}
