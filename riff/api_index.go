package riff

import (
	"fmt"
	"github.com/gimke/cart"
	"github.com/gimke/riff/common"
	"net/http"
)

type httpAPI struct{}

func (a *httpAPI) Index(r *cart.Router) {
	r.Route("/").GET(func(c *cart.Context) {
		c.Redirect(302, "/console/")
	})
	r.Route("/api", a.apiIndex)
}

func (a *httpAPI) apiIndex(r *cart.Router) {
	r.Route("").GET(a.version)
	r.Route("/version").GET(a.version)
	r.Route("/snap").GET(a.snap)
	r.Route("/nodes").GET(a.nodes)
	r.Route("/node/:name").GET(a.node)
	r.Route("/services").GET(a.services)
	r.Route("/service/:name", func(router *cart.Router) {
		router.Route("").GET(a.service)
		router.Route("/:command").GET(a.service)
		router.Route("/:command").POST(a.serviceCommand)

	})
	r.Route("/logs").GET(a.logs)
}

func (a httpAPI) version(c *cart.Context) {
	version := fmt.Sprintf("Cart version %s Riff version %s, build %s-%s", cart.Version, common.Version, common.GitBranch, common.GitSha)
	c.IndentedJSON(200, cart.H{
		"version": version,
	})
}

func (a httpAPI) snap(c *cart.Context) {
	c.IndentedJSON(200, cart.H{
		"snapShot": server.SnapShot,
	})
}

func (a httpAPI) nodes(c *cart.Context) {
	c.IndentedJSON(200, server.api.Nodes())
}

func (a httpAPI) node(c *cart.Context) {
	name, _ := c.Param("name")
	c.IndentedJSON(200, server.api.Node(name))
}

func (a httpAPI) services(c *cart.Context) {
	c.IndentedJSON(200, server.api.Services())
}

func (a httpAPI) service(c *cart.Context) {
	name, _ := c.Param("name")
	command, _ := c.Param("command")
	if command == "alive" {
		c.IndentedJSON(200, server.api.Service(name, false))
	} else {
		c.IndentedJSON(200, server.api.Service(name, true))
	}
}

func (a httpAPI) serviceCommand(c *cart.Context) {
	name, _ := c.Param("name")
	command, _ := c.Param("command")
	var err error
	code := 200
	switch command {
	case "start":
		code = server.api.Start(name)
		break
	case "stop":
		code = server.api.Stop(name)
		break
	case "restart":
		code = server.api.Restart(name)
		break
	default:
		code = 400
	}

	switch code {
	case 200:
		c.IndentedJSON(code, cart.H{
			"status": code,
		})
	case 201:
		c.IndentedJSON(code, cart.H{
			"status": code,
		})
	case 404:
		err = fmt.Errorf("not found")
		c.IndentedJSON(code, cart.H{
			"status": code,
			"error":  err.Error(),
		})
	case 400:
		err = fmt.Errorf("command missing")
		c.IndentedJSON(code, cart.H{
			"status": code,
			"error":  err.Error(),
		})
	default:
		err = fmt.Errorf("error")
		c.IndentedJSON(code, cart.H{
			"status": code,
			"error":  err.Error(),
		})
	}

}

type httpLogHandler struct {
	logCh chan string
}

func (h *httpLogHandler) HandleLog(log string) {
	// Do a non-blocking send
	select {
	case h.logCh <- log:
	}
}

func (a httpAPI) logs(c *cart.Context) {
	resp := c.Response
	clientGone := resp.(http.CloseNotifier).CloseNotify()

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.Header().Set("Connection", "Keep-Alive")
	resp.Header().Set("Transfer-Encoding", "chunked")
	resp.Header().Set("X-Content-Type-Options", "nosniff")

	handler := &httpLogHandler{
		logCh: make(chan string, 512),
	}
	server.logWriter.RegisterHandler(handler)
	defer server.logWriter.DeregisterHandler(handler)

	flusher, ok := resp.(http.Flusher)
	if !ok {
		server.Logger.Println("Streaming not supported")
	}
	for {
		select {
		case <-clientGone:
			return
		case logs := <-handler.logCh:
			fmt.Fprintln(resp, logs)
			flusher.Flush()
		}
	}
}
