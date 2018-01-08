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
	s.logger.Printf(infoRpcPrefix+"start to accept http conn: %v", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.Shutdown()
		s.logger.Printf(errorRpcPrefix+"start http server error: %s", err)
	}
}
func (s *Server) listenRpc() {
	s.logger.Printf(infoRpcPrefix+"start to accept rpc conn: %v", s.Listener.Addr())
	s.logger.Printf(infoRpcPrefix+"riff snapshot now is: %s", s.SnapShot)
	for {
		// Accept a connection
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.shutdown {
				return
			}
			s.logger.Printf(errorRpcPrefix+"failed to accept RPC conn: %v", err)
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
				s.logger.Printf(infoRpcPrefix+"end of %s", conn.RemoteAddr().String())
			} else {
				s.logger.Printf(errorRpcPrefix+"%v %s", err, conn.RemoteAddr().String())
			}
			return
		}
	}
}

func (s *Server) print() {
	fmt.Printf(`
    Riff running!

        Node Id:  %v
           Name:  %v
             DC:  %v
   HTTP Address:  %v
    RPC Address:  %v

`, s.Self.Id, s.Self.Name, s.Self.DataCenter, s.httpServer.Addr, s.Listener.Addr())
}
