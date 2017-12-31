package riff

import (
	"github.com/gimke/riff/common"
	"io"
	"log"
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
}

func (s *Server) listen() {
	for {
		// Accept a connection
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.shutdown {
				return
			}
			log.Printf("[ERR] riff.rpc: failed to accept RPC conn: %v", err)
			continue
		}
		go s.handleConn(conn)
		//metrics.IncrCounter([]string{"riff", "rpc", "accept_conn"}, 1)
		//metrics.IncrCounter([]string{"rpc", "accept_conn"}, 1)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	codec := common.NewGobServerCodec(conn)
	for {
		select {
		case <-s.shutdownCh:
			return
		default:
		}
		if err := s.rpcServer.ServeRequest(codec); err != nil {
			if err == io.EOF {
				log.Printf("[INFO] riff.rpc: end of %s", conn.RemoteAddr().String())
			} else {
				log.Printf("[ERR] riff.rpc: %v %s", err, conn.RemoteAddr().String())
			}
			return
		}
	}
}
