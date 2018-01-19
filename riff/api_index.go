package riff

import (
	"fmt"
	"github.com/gimke/cart"
	"github.com/gimke/riff/common"
	"net/http"
)

type httpAPI struct {}

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
	r.Route("/service/:name").GET(a.service)
	r.Route("/logs").GET(a.logs)
}

func (a httpAPI) version(c *cart.Context) {
	version := fmt.Sprintf("Cart version %s Riff version %s, build %s-%s", cart.Version, common.Version, common.GitBranch, common.GitSha)
	c.IndentedJSON(200, cart.H{
		"Version": version,
	})
}

func (a httpAPI) snap(c *cart.Context) {
	c.IndentedJSON(200, cart.H{
		"SnapShot": server.SnapShot,
	})
}

func (a httpAPI) nodes(c *cart.Context) {
	c.IndentedJSON(200, server.api.Nodes())
}

func (a httpAPI) services(c *cart.Context) {
	c.IndentedJSON(200, server.ServicesSlice())
}

func (a httpAPI) service(c *cart.Context) {
	name, _ := c.Param("name")
	c.IndentedJSON(200, server.GetService(name))
}

func (a httpAPI) node(c *cart.Context) {
	name, _ := c.Param("name")
	c.IndentedJSON(200, server.GetNode(name))
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
