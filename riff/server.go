package riff

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Server struct {
	Listener     net.Listener
	rpcServer    *rpc.Server
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

func NewServer() (*Server, error) {
	shutdownCh := make(chan struct{})

	s := &Server{
		rpcServer:  rpc.NewServer(),
		shutdownCh: shutdownCh,
	}
	if err := s.setupRPC(); err != nil {
		s.Shutdown()
		return nil, fmt.Errorf("Failed to start RPC layer: %v", err)
	}
	go s.listen()
	return s, nil
}

func (s *Server) setupRPC() error {
	for _, fn := range endpoints {
		s.rpcServer.Register(fn(s))
	}
	//s.rpcServer.Register(&Status{server:s})
	addr, err := net.ResolveTCPAddr("", ":8530")
	if err != nil {
		return err
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
