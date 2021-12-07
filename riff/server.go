package riff

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/teatak/cart"
	"github.com/teatak/riff/api"
	"github.com/teatak/riff/common"
	"github.com/teatak/riff/schema"
)

const errorServerPrefix = "[ERR]  riff.server: "
const infoServerPrefix = "[INFO] riff.server: "
const errorServicePrefix = "[ERR]  riff.service: "
const infoServicePrefix = "[INFO] riff.service: "

var server *Server

type Server struct {
	Listener   net.Listener
	Logger     *log.Logger
	Self       *Node
	ShutdownCh chan struct{}
	SnapShot   string
	Nodes
	api          *API
	httpServer   *http.Server
	rpcServer    *rpc.Server
	logWriter    *LogWriter
	watch        *Watch
	config       *Config
	shutdown     bool
	shutdownLock sync.Mutex
}

func NewServer(config *Config) (*Server, error) {

	shutdownCh := make(chan struct{})

	server = &Server{
		logWriter:  NewLogWriter(255),
		watch:      NewWatch(),
		rpcServer:  rpc.NewServer(),
		api:        &API{},
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
	ip, _, err := net.ParseCIDR(s.config.IP)
	if err != nil {
		s.Logger.Printf(errorServicePrefix+" %v\n", err)
	}
	self := &Node{
		Name:        s.config.Name,
		IP:          ip.String(),
		RpcPort:     s.config.Ports.Rpc,
		HttpPort:    s.config.Ports.Http,
		DataCenter:  s.config.DataCenter,
		IsSelf:      true,
		State:       api.StateAlive,
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

var checkHtpasswd = false
var basicAuth = ""

func (s *Server) auth(c *cart.Context, next cart.Next) {
	//get htpasswd

	if !checkHtpasswd {
		//check
		file := common.BinDir + "/config/.htpasswd"
		if common.IsExist(file) {
			content, _ := ioutil.ReadFile(file)
			basicAuth = string(content)
		}
		checkHtpasswd = true
	}

	if checkHtpasswd && basicAuth != "" {
		//需要验证
		r := c.Request
		user, pass, ok := r.BasicAuth()
		if !ok {
			//check user pass
			arr := strings.Split(basicAuth, ":")
			cuser := arr[0]
			cpass := strings.ReplaceAll(arr[1], "{SHA}", "")
			if user != cuser || common.ShaPass(pass) != cpass {
				c.Response.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
				c.Response.WriteHeader(401)
				c.JSON(401, cart.H{
					"error": "Unauthorized",
				})
				return
			}
		}
	}

	//if s.config.Auth != nil {
	//	r := c.Request
	//	user, pass, ok := r.BasicAuth()
	//	if !ok || user != s.config.Auth.UserName || pass != s.config.Auth.Password {
	//		c.Response.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
	//		c.Response.WriteHeader(401)
	//		c.JSON(401, cart.H{
	//			"error": "Unauthorized",
	//		})
	//		return
	//	}
	//}
	next()
}

func (s *Server) setupCart() error {
	cart.SetMode(cart.ReleaseMode)
	//http.Handle("/", http.FileServer(assetFS()))
	//http.ListenAndServe(s.config.Addresses.Http + ":" + strconv.Itoa(s.config.Ports.Http),nil)

	r := cart.New()
	r.Use("/", s.auth, s.httpLogger(), cart.RecoveryRender(cart.DefaultErrorWriter))
	if common.IsRelease() {
		r.Use("/favicon.ico", func(c *cart.Context, next cart.Next) {
			b, err := assetFS().Asset("static/images/favicon.ico")
			if err != nil {
				s.Logger.Printf(errorServerPrefix+"error: %v\n", err)
				next()
			} else {
				c.Response.WriteHeader(200)
				c.Response.Write(b)
			}
		})
		r.Use("/console/*file", func(c *cart.Context, next cart.Next) {
			b, err := assetFS().Asset("static/dist/console.html")
			if err != nil {
				s.Logger.Printf(errorServerPrefix+"error: %v\n", err)
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
	} else {
		r.Use("/favicon.ico", cart.File("./static/images/favicon.ico"))
		r.Use("/console/*file", cart.File("./static/dist/console.html"))
		r.Use("/static/*file", cart.Static("./static", false))
	}
	h := &Http{Schema: graphql.MustParseSchema(schema.String(), &Resolver{})}
	r.Route("/", h.Index)
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

	s.Logger.Printf(infoServerPrefix+"%s leave", s.Self.Name)
	s.fanoutLeave()

	s.shutdown = true
	close(s.ShutdownCh)
	if s.Listener != nil {
		s.Listener.Close()
	}
	return nil
}
