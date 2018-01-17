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
const errorNodePrefix = "[ERR]  riff.node: "
const infoNodePrefix = "[INFO] riff.node: "
const errorServicePrefix = "[ERR]  riff.service: "
const infoServicePrefix = "[INFO] riff.service: "

var server *Server

type Server struct {
	Listener   net.Listener
	httpServer *http.Server
	rpcServer  *rpc.Server
	Logger     *log.Logger
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

	server = &Server{
		logWriter:  NewLogWriter(512),
		rpcServer:  rpc.NewServer(),
		config:     config,
		ShutdownCh: shutdownCh,
	}

	logOutput := io.MultiWriter(os.Stderr, server.logWriter)
	server.Logger = log.New(logOutput, "", log.LstdFlags|log.Lmicroseconds)

	if err := server.setupServer(); err != nil {
		server.Shutdown()
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}

	if err := server.setupRPC(); err != nil {
		server.Shutdown()
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	if err := server.setupCart(); err != nil {
		server.Shutdown()
		return nil, fmt.Errorf(errorServerPrefix+"%v", err)
	}
	server.print()
	server.initServices()
	go server.handleServices()  //handle service
	go server.listenRpc()       //listen rpc
	go server.listenHttp()      //listen http
	go server.fanoutNodes()     //fanout state
	go server.fanoutDeadNodes() //fanout dead state
	return server, nil
}
func (s *Server) setupServer() error {
	self := &Node{
		Name:        s.config.Name,
		IP:          s.config.Addresses.Rpc,
		Port:        s.config.Ports.Rpc,
		DataCenter:  s.config.DataCenter,
		IsSelf:      true,
		State:       stateAlive,
		StateChange: time.Now(),
		Services:    make(map[string]*Service),
	}
	s.Self = self
	s.AddNode(self)
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

		s.Logger.Printf("[INFO] cart: status:%d latency:%v ip:%s method:%s path:%s\n",
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
			s.Logger.Printf(errorRpcPrefix+"error: %v\n", err)
			next()
		} else {
			c.Response.WriteHeader(200)
			c.Response.Write(b)
		}
	})
	//r.Use("/console/*file", func(c *cart.Context, next cart.Next) {
	//	b, err := assetFS().Asset("static/dist/console.html")
	//	if err != nil {
	//		s.Logger.Printf(errorRpcPrefix+"error: %v\n", err)
	//		next()
	//	} else {
	//		c.Response.WriteHeader(200)
	//		c.Response.Write(b)
	//	}
	//})
	//r.Use("/static/*file", func(c *cart.Context, next cart.Next) {
	//	http.StripPrefix("/static/", http.FileServer(assetFS())).ServeHTTP(c.Response, c.Request)
	//	if c.Response.Status() == 404 {
	//		c.Response.WriteHeader(200) //reset status
	//		next()
	//	}
	//})
	//debug
	r.Use("/console/*file", cart.File("../static/dist/console.html"))
	r.Use("/static/*file", cart.Static("../static", false))
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

	s.Logger.Printf(infoRpcPrefix+"%s leave", s.Self.Name)
	s.fanoutLeave()

	s.shutdown = true
	close(s.ShutdownCh)
	if s.Listener != nil {
		s.Listener.Close()
	}
	return nil
}
