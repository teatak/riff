package riff

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

const errorServerPrefix = "riff.server error: "
const errorRpcPrefix = "[ERR]  riff.rpc error: "
const infoRpcPrefix = "[INFO] riff.rpc: "

type Server struct {
	Listener   net.Listener
	rpcServer  *rpc.Server
	Name       string
	DataCenter string
	Nodes
	Services
	SnapShort    string
	config       *Config
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

func NewServer(config *Config) (*Server, error) {
	shutdownCh := make(chan struct{})

	s := &Server{
		rpcServer:  rpc.NewServer(),
		config:     config,
		shutdownCh: shutdownCh,
	}

	s.setupServer()
	if err := s.setupRPC(); err != nil {
		s.Shutdown()
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	go s.listen()
	return s, nil
}
func (s *Server) setupServer() error {
	self := &Node{
		Name:       s.config.Name,
		IP:         s.config.IP,
		Port:       s.config.Port,
		DataCenter: s.config.DataCenter,
		State:      stateAlive,
	}
	s.Name = s.config.Name
	s.DataCenter = s.config.DataCenter
	s.Nodes = make(map[string]*Node)
	s.Services = make(map[string]*Service)
	s.AddNode(self)
	return nil
}

func (s *Server) setupRPC() error {
	for _, fn := range endpoints {
		s.rpcServer.Register(fn(s))
	}
	addr := &net.TCPAddr{
		IP:   s.config.IP,
		Port: s.config.Port,
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	s.Listener = ln
	return nil
}

func (s *Server) Shutdown() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}

	s.shutdown = true
	close(s.shutdownCh)
	if s.Listener != nil {
		s.Listener.Close()
	}

	return nil
}
