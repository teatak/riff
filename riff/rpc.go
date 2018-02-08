package riff

import (
	"fmt"
	"github.com/gimke/riff/api"
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
	registerEndpoint(func(s *Server) interface{} { return &Riff{} })
	registerEndpoint(func(s *Server) interface{} { return &Query{} })
	registerEndpoint(func(s *Server) interface{} { return &Mutation{} })
}

func (s *Server) listenHttp() {
	s.Logger.Printf(infoServerPrefix+"start to accept http conn: %v\n", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.Shutdown()
		s.Logger.Printf(errorServerPrefix+"start http server error: %s\n", err)
	}
}
func (s *Server) listenRpc() {
	s.Logger.Printf(infoServerPrefix+"start to accept rpc conn: %v\n", s.Listener.Addr())
	for {
		// Accept a connection
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.shutdown {
				return
			}
			s.Logger.Printf(errorServerPrefix+"failed to accept RPC conn: %v\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1)
	if _, err := conn.Read(buf); err != nil {
		if err != io.EOF {
			s.Logger.Printf(errorServerPrefix+"failed to read conn: %v\n", err)
		}
		conn.Close()
		return
	}
	typ := api.RPCType(buf[0])
	switch typ {
	case api.RPCRiff:
		s.handleRpc(conn)
	case api.RPCLog:
		s.handleLog(conn)
	}
}

func (s *Server) handleRpc(conn net.Conn) {
	codec := api.NewGobServerCodec(conn)
	for {
		select {
		case <-s.ShutdownCh:
			return
		default:
		}
		if err := s.rpcServer.ServeRequest(codec); err != nil {
			if err == io.EOF {
			} else {
				s.Logger.Printf(errorServerPrefix+"%v %s\n", err, conn.RemoteAddr().String())
			}
			return
		}
	}
}

func (s *Server) handleLog(conn net.Conn) {
	codec := api.NewGobServerCodec(conn)
	for {
		select {
		case <-s.ShutdownCh:
			return
		default:
		}
		if err := s.rpcServer.ServeRequest(codec); err != nil {
			if err == io.EOF {
			} else {
				s.Logger.Printf(errorServerPrefix+"%v %s\n", err, conn.RemoteAddr().String())
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
