package riff

import (
	"fmt"
	"github.com/gimke/cart"
	"github.com/gimke/riff/common"
	"net/http"
)

type Api struct {
	server *Server
}

func (a *Api) Index(r *cart.Router) {
	r.Route("/").GET(func(c *cart.Context) {
		c.Redirect(302, "/console/")
	})
	r.Route("/api", a.apiIndex)
}

func (a *Api) apiIndex(r *cart.Router) {
	r.Route("").GET(a.version)
	r.Route("/version").GET(a.version)
	r.Route("/snap").GET(a.snap)
	r.Route("/nodes").GET(a.nodes)
	r.Route("/node/:name").GET(a.node)
	r.Route("/services").GET(a.services)
	r.Route("/logs").GET(a.logs)
}

func (a Api) version(c *cart.Context) {
	version := fmt.Sprintf("Cart version %s Riff version %s, build %s-%s", cart.Version, common.Version, common.GitBranch, common.GitSha)
	c.IndentedJSON(200, cart.H{
		"Version": version,
	})
}

func (a Api) snap(c *cart.Context) {
	c.IndentedJSON(200, cart.H{
		"SnapShot": a.server.SnapShot,
	})
}

func (a Api) nodes(c *cart.Context) {
	c.IndentedJSON(200, a.server.Slice())
}

func (a Api) services(c *cart.Context) {
	c.IndentedJSON(200, a.server.ServicesSlice())
}

func (a Api) node(c *cart.Context) {
	name, _ := c.Param("name")
	c.IndentedJSON(200, a.server.GetNode(name))
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

func (a Api) logs(c *cart.Context) {
	resp := c.Response
	clientGone := resp.(http.CloseNotifier).CloseNotify()

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.Header().Set("Connection", "Keep-Alive")
	resp.Header().Set("Transfer-Encoding", "chunked")
	resp.Header().Set("X-Content-Type-Options", "nosniff")

	handler := &httpLogHandler{
		logCh: make(chan string, 512),
	}
	a.server.logWriter.RegisterHandler(handler)
	defer a.server.logWriter.DeregisterHandler(handler)

	flusher, ok := resp.(http.Flusher)
	if !ok {
		a.server.Logger.Println("Streaming not supported")
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
