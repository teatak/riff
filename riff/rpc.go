package riff

import (
	"fmt"
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
	registerEndpoint(func(s *Server) interface{} { return &Query{s} })
}

func (s *Server) listen() {
	s.print()
	log.Printf(infoRpcPrefix+"start to accept rpc conn: %v", s.Listener.Addr())
	log.Printf(infoRpcPrefix+"snapshot now is:%s", s.SnapShot)
	for {
		// Accept a connection
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.shutdown {
				return
			}
			log.Printf(errorRpcPrefix+"failed to accept RPC conn: %v", err)
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
		case <-s.shutdownCh:
			return
		default:
		}
		if err := s.rpcServer.ServeRequest(codec); err != nil {
			if err == io.EOF {
				log.Printf(infoRpcPrefix+"end of %s", conn.RemoteAddr().String())
			} else {
				log.Printf(errorRpcPrefix+"%v %s", err, conn.RemoteAddr().String())
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
    RPC Address:  %v

`, s.Id, s.Name, s.DataCenter, s.Listener.Addr())
}
