package riff

import (
	"fmt"
	"github.com/gimke/cart"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
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
	logger     *log.Logger
	logWriter  *LogWriter
	Self       *Node
	Nodes
	SnapShot     string
	config       *Config
	shutdown     bool
	ShutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

func NewServer(config *Config) (*Server, error) {

	shutdownCh := make(chan struct{})

	s := &Server{
		logWriter:  NewLogWriter(512),
		rpcServer:  rpc.NewServer(),
		config:     config,
		ShutdownCh: shutdownCh,
	}

	logOutput := io.MultiWriter(os.Stderr, s.logWriter)
	s.logger = log.New(logOutput, "", log.LstdFlags)

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
	s.Self = self
	s.Nodes = make(map[string]*Node)
	s.AddNode(self)
	s.Shutter() //make snap sort
	return nil
}

func (s *Server) httpLogger() cart.Handler {
	return func(c *cart.Context, next cart.Next) {
		start := time.Now()
		path := c.Request.URL.Path
		next()
		end := time.Now()
		latency := end.Sub(start)
		method := c.Request.Method
		clientIP := c.ClientIP()
		statusCode := c.Response.Status()

		s.logger.Printf("[INFO] cart: status:%d latency:%v ip:%s method:%s path:%s\n",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

func (s *Server) setupCart() error {
	cart.SetMode(cart.ReleaseMode)
	//http.Handle("/", http.FileServer(assetFS()))
	//http.ListenAndServe(s.config.Addresses.Http + ":" + strconv.Itoa(s.config.Ports.Http),nil)
	//
	r := cart.New()
	r.Use("/", s.httpLogger(), cart.RecoveryRender(cart.DefaultErrorWriter))
	r.Use("/favicon.ico", func(c *cart.Context, next cart.Next) {
		b, err := assetFS().Asset("static/images/favicon.ico")
		if err != nil {
			s.logger.Printf(errorRpcPrefix+"error: %v\n", err)
			next()
		} else {
			c.Response.WriteHeader(200)
			c.Response.Write(b)
		}
	})
	r.Use("/console/*file", func(c *cart.Context, next cart.Next) {
		b, err := assetFS().Asset("static/index.html")
		if err != nil {
			s.logger.Printf(errorRpcPrefix+"error: %v\n", err)
			next()
		} else {
			c.Response.WriteHeader(200)
			c.Response.Write(b)
		}
	})
	r.Use("/static/*file", func(c *cart.Context, next cart.Next) {
		http.StripPrefix("/static/", http.FileServer(assetFS())).ServeHTTP(c.Response, c.Request)
		if c.Response.Status() == 404 {
			c.Response.WriteHeader(200) //reset status
			next()
		}
	})
	a := Api{
		server: s,
	}
	r.Route("/", a.Index)
	s.httpServer = r.ServerKeepAlive(s.config.Addresses.Http + ":" + strconv.Itoa(s.config.Ports.Http))
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
	close(s.ShutdownCh)
	if s.Listener != nil {
		s.Listener.Close()
	}

	return nil
}
