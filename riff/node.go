package riff

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/gimke/riff/api"
	"github.com/gimke/riff/common"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

func (s *Server) Keys() []string {
	var keys = make([]string, 0, 0)
	s.nodes.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	/*sort.Slice(keys, func(i, j int) bool {
		var validID = regexp.MustCompile(`^(\S+)([0-9]+)$`)
		if validID.MatchString(keys[i]) && validID.MatchString(keys[j]) {
			arri := validID.FindAllStringSubmatch(keys[i],-1)
			arrj := validID.FindAllStringSubmatch(keys[j],-1)
			if arri[0][1] == arrj[0][1] {
				inti,_ := strconv.Atoi(arri[0][2])
				intj,_ := strconv.Atoi(arrj[0][2])
				return inti <  intj
			} else {
				return  arri[0][1] < arrj[0][1]
			}
		} else {
			return keys[i] < keys[j]
		}
	})*/
	sort.Strings(keys)
	return keys
}

func (s *Server) Slice() []*Node {
	keys := s.Keys()
	nodes := make([]*Node, 0)
	for _, key := range keys {
		if n := s.GetNode(key); n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (s *Server) ServicesSlice() []string {
	keys := s.Keys()
	helper := make(map[string]string, 0)
	services := []string{}
	for _, key := range keys {
		if n := s.GetNode(key); n != nil {
			for name, _ := range n.Services {
				if _, ok := helper[name]; !ok {
					helper[name] = name
					services = append(services, name)
				}
			}
		}
	}
	return services
}

func (s *Server) GetService(findName string) interface{} {
	keys := s.Keys()
	var service map[string]interface{}
	nodes := []string{}
	for _, key := range keys {
		if n := s.GetNode(key); n != nil {
			for name, s := range n.Services {
				if name == findName {
					if service == nil {
						service = map[string]interface{}{
							"Name": name,
						}
					}
					if n.State == api.StateAlive {
						nodes = append(nodes, s.Address())
					}
				}
			}
		}
	}
	if service != nil {
		service["Nodes"] = nodes
	}
	return service
}

func (s *Server) Range(f func(string, *Node) bool) {
	s.nodes.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Node))
	})
}

func (s *Server) GetNode(key string) *Node {
	if n, _ := s.nodes.Load(key); n != nil {
		return n.(*Node)
	}
	return nil
}

func (s *Server) AddNode(n *Node) {
	n.Shutter()
	s.nodes.Store(n.Name, n)

	//watch
	s.watch.Dispatch(WatchParam{
		Name:      n.Name,
		WatchType: NodeChanged,
	})
}
func (s *Server) AddService(service *Service) {
	s.Self.Services[service.Name] = service
}
func (s *Server) DeleteNode(key string) {
	s.nodes.Delete(key)

	//watch
	s.watch.Dispatch(WatchParam{
		Name:      key,
		WatchType: NodeChanged,
	})
}

func (s *Server) RemoveNodeDelay(n *Node) {
	if n.timeoutFn == nil {
		n.timeoutFn = func() {
			//delete this node
			if n.State == api.StateDead {
				s.Logger.Printf(infoServerPrefix+"remove dead node %s\n", n.Name)
				s.DeleteNode(n.Name)
				//clear fn
				n.timeoutFn = nil
			}
		}
		timeout := 5 * 60 * time.Second
		n.timer = time.AfterFunc(timeout, n.timeoutFn)
	}
}

// set node state version inc and shutter node return version
func (s *Server) SetStateOnly(n *Node, state api.StateType) uint64 {
	n.State = state
	n.Shutter()
	return n.VersionGet()
}

// set node state version inc and shutter node return version++
func (s *Server) SetState(n *Node, state api.StateType) uint64 {
	n.State = state
	v := n.VersionInc()
	n.Shutter()
	return v
}

// SetState and make a snapsort return version ++
func (s *Server) SetStateWithShutter(n *Node, state api.StateType) uint64 {
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

type Node struct {
	Name        string
	DataCenter  string
	IP          string
	RpcPort     int
	HttpPort    int
	Version     uint64
	State       api.StateType // Current state
	StateChange time.Time     // Time last state change happened
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
	return net.JoinHostPort(n.IP, strconv.Itoa(int(n.RpcPort)))
}

func (n *Node) String() string {
	buff := bytes.NewBuffer(nil)
	io.WriteString(buff, n.Name+":"+strconv.FormatUint(n.Version, 10)+":{")
	//write service name and version
	keys := n.Services.Keys()
	for i, sk := range keys {
		s := n.Services[sk]
		io.WriteString(buff, s.Name+":{"+s.Address()+","+s.State.Name()+","+s.StatusContent+"}")
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

func (n *Node) LoadServices() {
	if n.Services == nil {
		n.Services = make(map[string]*Service)
	}
	files, err := ioutil.ReadDir(common.BinDir + "/config")
	if err == nil {
		for _, file := range files {
			basename := file.Name()
			if basename == "riff.yml" {
				continue
			}
			if strings.HasPrefix(basename, ".") {
				continue
			}
			name := strings.TrimSuffix(basename, filepath.Ext(basename))
			s := n.LoadService(name)
			if s != nil {
				n.Services[s.Name] = s
			}
		}
	}
	n.Shutter()
}

func (n *Node) LoadService(name string) *Service {
	file := common.BinDir + "/config/" + name + ".yml"
	if !common.IsExist(file) {
		return nil
	}
	content, _ := ioutil.ReadFile(file)
	var s = &Service{}
	var c = &ServiceConfig{}
	err := yaml.Unmarshal(content, &c)

	if err != nil {
		server.Logger.Printf(errorServicePrefix+"%s config file error: %v", name, err)
		return nil
	}
	s.Config = string(content)
	s.ServiceConfig = c
	s.StateChange = time.Now()
	s.runAtLoad()
	s.checkState()
	return s
}
