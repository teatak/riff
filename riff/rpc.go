package riff

import (
	"fmt"
	"github.com/gimke/riff/common"
	"io"
	"net"
)

type factory func(s *Server) interface{}

// endpoints is a list of registered RPC endpoint factories.
var endpoints []factory

// registerEndpoint registers a new RPC endpoint factory.
func registerEndpoint(fn factory) {
	endpoints = append(endpoints, fn)
}

func init() {
	registerEndpoint(func(s *Server) interface{} { return &Status{s} })
	registerEndpoint(func(s *Server) interface{} { return &Query{s} })
	registerEndpoint(func(s *Server) interface{} { return &Riff{s} })
}

func (s *Server) listenHttp() {
	s.Logger.Printf(infoRpcPrefix+"start to accept http conn: %v\n", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.Shutdown()
		s.Logger.Printf(errorRpcPrefix+"start http server error: %s\n", err)
	}
}
func (s *Server) listenRpc() {
	s.Logger.Printf(infoRpcPrefix+"start to accept rpc conn: %v\n", s.Listener.Addr())
	for {
		// Accept a connection
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.shutdown {
				return
			}
			s.Logger.Printf(errorRpcPrefix+"failed to accept RPC conn: %v\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	codec := common.NewGobServerCodec(conn)
	for {
		select {
		case <-s.ShutdownCh:
			return
		default:
		}
		if err := s.rpcServer.ServeRequest(codec); err != nil {
			if err == io.EOF {
			} else {
				s.Logger.Printf(errorRpcPrefix+"%v %s\n", err, conn.RemoteAddr().String())
			}
			return
		}
	}
}

func (s *Server) print() {
	fmt.Printf(`
    Riff running!

           Name:  %v
             DC:  %v
   HTTP Address:  %v
    RPC Address:  %v

`, s.Self.Name, s.Self.DataCenter, s.httpServer.Addr, s.Listener.Addr())
}
