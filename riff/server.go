package riff

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

const errorServerPrefix = "riff.server error: "
const errorRpcPrefix = "[ERR]  riff.rpc: "
const infoRpcPrefix = "[INFO] riff.rpc: "

type Server struct {
	Listener  net.Listener
	rpcServer *rpc.Server
	Id        string
	Nodes
	Services
	SnapShot     string
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
		Id:          s.config.Id,
		Name:        s.config.Name,
		IP:          s.config.Addresses.Rpc,
		Port:        s.config.Ports.Rpc,
		DataCenter:  s.config.DataCenter,
		State:       stateAlive,
		StateChange: time.Now(),
	}
	s.Id = s.config.Id
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
		IP:   net.ParseIP(s.config.Addresses.Rpc),
		Port: s.config.Ports.Rpc,
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
