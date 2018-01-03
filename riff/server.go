package riff

import (
	"fmt"
	"github.com/gimke/cart"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

const errorServerPrefix = "riff.server error: "
const errorRpcPrefix = "[ERR]  riff.rpc: "
const infoRpcPrefix = "[INFO] riff.rpc: "

type Server struct {
	Listener   net.Listener
	httpServer *http.Server
	rpcServer  *rpc.Server
	Id         string
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

	if err := s.setupServer(); err != nil {
		s.Shutdown()
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}

	if err := s.setupRPC(); err != nil {
		s.Shutdown()
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	if err := s.setupCart(); err != nil {
		s.Shutdown()
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	s.print()
	go s.listenRpc()
	go s.listenHttp()
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

func assetServer() cart.Handler {
	return func(c *cart.Context, next cart.Next) {
		http.StripPrefix("/console/", http.FileServer(assetFS())).ServeHTTP(c.Response, c.Request)
	}
}

func (s *Server) setupCart() error {
	cart.SetMode(cart.ReleaseMode)
	//http.Handle("/", http.FileServer(assetFS()))
	//http.ListenAndServe(s.config.Addresses.Http + ":" + strconv.Itoa(s.config.Ports.Http),nil)
	//
	r := cart.New()
	r.Use("/", Logger(), cart.RecoveryRender(cart.DefaultErrorWriter))
	r.Use("/console/*file", assetServer())
	a := Api{
		server: s,
	}
	r.Route("/", a.Index)
	s.httpServer = r.Server(s.config.Addresses.Http + ":" + strconv.Itoa(s.config.Ports.Http))
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
