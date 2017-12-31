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
	Listener     net.Listener
	rpcServer    *rpc.Server
	riff         *Riff
	config       *Config
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

func NewServer(config *Config) (*Server, error) {
	shutdownCh := make(chan struct{})

	riff, err := Create(config.Name)
	if err != nil {
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	if err != nil {
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	self := &Node{
		Name:  config.Name,
		IP:    config.IP,
		Port:  config.Port,
		State: stateAlive,
	}
	riff.AddNode(self)

	s := &Server{
		rpcServer:  rpc.NewServer(),
		config:     config,
		shutdownCh: shutdownCh,
		riff:       riff,
	}

	if err := s.setupRPC(); err != nil {
		s.Shutdown()
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	go s.listen()
	return s, nil
}

func (s *Server) setupRPC() error {
	for _, fn := range endpoints {
		s.rpcServer.Register(fn(s))
	}
	//s.rpcServer.Register(&Status{server:s})
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
