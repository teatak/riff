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
	"sync/atomic"
	"time"
)

type Nodes struct {
	nodes sync.Map
}

type Digest struct {
	Name     string
	SnapShot string
}

var removeFirst = 0

func (s *Server) Shutter() {
	h := sha1.New()
	io.WriteString(h, s.String())
	s.SnapShot = fmt.Sprintf("%x", h.Sum(nil))
	if removeFirst != 0 {
		s.Logger.Printf(infoRpcPrefix+"riff snapshot now is: %s\n", s.SnapShot)
	} else {
		removeFirst++
	}

}

func (s *Server) Keys() []string {
	var keys = make([]string, 0, 0)
	s.nodes.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	sort.Strings(keys)
	return keys
}

func (s *Server) Slice() []*Node {
	keys := s.Keys()
	nodes := make([]*Node, 0)
	for _, key := range keys {
		if n, _ := s.nodes.Load(key); n != nil {
			nodes = append(nodes, n.(*Node))
		}
	}
	return nodes
}
func (s *Server) Range(f func(string, *Node) bool) {
	s.nodes.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Node))
	})
}
func (ns *Server) GetNode(key string) *Node {
	if n, _ := ns.nodes.Load(key); n != nil {
		return n.(*Node)
	}
	return nil
}

func (s *Server) AddNode(n *Node) {
	n.Shutter()
	s.nodes.Store(n.Name, n)
}

func (s *Server) DeleteNode(key string) {
	s.nodes.Delete(key)
}

func (s *Server) RemoveNodeDelay(n *Node) {
	if n.timeoutFn == nil {
		n.timeoutFn = func() {
			//delete this node
			if n.State == stateDead {
				s.Logger.Printf(infoNodePrefix+"remove dead node %s\n", n.Name)
				s.DeleteNode(n.Name)
				//clear fn
				n.timeoutFn = nil
			}
		}
		timeout := 10 * time.Second
		n.timer = time.AfterFunc(timeout, n.timeoutFn)
	}
}

// set node state version inc and shutter node return version
func (s *Server) SetStateOnly(n *Node, state stateType) uint64 {
	n.State = state
	n.Shutter()
	return n.VersionGet()
}

// set node state version inc and shutter node return version++
func (s *Server) SetState(n *Node, state stateType) uint64 {
	n.State = state
	v := n.VersionInc()
	n.Shutter()
	return v
}

// SetState and make a snapsort return version ++
func (s *Server) SetStateWithShutter(n *Node, state stateType) uint64 {
	v := s.SetState(n, state)
	s.Shutter()
	return v
}

// get random nodes
func (s *Server) randomNodes(fanout int, filterFn func(*Node) bool) []*Node {
	nodes := s.Keys()
	n := len(nodes)
	RNodes := make([]*Node, 0, fanout)
OUTER:
	for i := 0; i < 3*n && len(RNodes) < fanout; i++ {
		idx := common.RandomNumber(n)
		n := s.GetNode(nodes[idx])
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

func (n *Node) VersionGet() uint64 {
	return atomic.LoadUint64(&n.Version)
}

func (n *Node) VersionInc() uint64 {
	return atomic.AddUint64(&n.Version, 1)
}

// Witness is called to update our local clock if necessary after
// witnessing a clock value received from another process
func (n *Node) VersionSet(v uint64) {
WITNESS:
	// If the other value is old, we do not need to do anything
	cur := atomic.LoadUint64(&n.Version)
	other := uint64(v)
	if other < cur {
		return
	}

	// Ensure that our local clock is at least one ahead.
	if !atomic.CompareAndSwapUint64(&n.Version, cur, other) {
		// The CAS failed, so we just retry. Eventually our CAS should
		// succeed or a future witness will pass us by and our witness
		// will end.
		goto WITNESS
	}
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

//func (n *Node) Suspect() {
//	n.State = stateSuspect
//	n.VersionInc()
//	n.Shutter()
//}
//
//func (n *Node) Alive() {
//	n.State = stateAlive
//	n.VersionInc()
//	n.Shutter()
//}

func (n *Node) AddService(s *Service) {
	if n.Services == nil {
		n.Services = make(map[string]*Service)
	}
	n.Services[s.Name] = s
	n.Shutter()
}
