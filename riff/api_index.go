package riff

import (
	"fmt"
	"github.com/gimke/cart"
	"time"
	"net/http"
	"log"
	"github.com/gimke/riff/common"
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
}

func (a Api) version(c *cart.Context) {
	version := fmt.Sprintf("Cart version %s Riff version %s, build %s-%s", cart.Version, common.Version, common.GitBranch, common.GitSha)
	//c.JSON(200, cart.H{
	//	"Version": version,
	//})
	resp := c.Response
	resp.Header().Set("Content-Type", "text/event-stream")
	resp.Header().Set("Cache-Control", "no-cache")
	resp.Header().Set("Connection", "keep-alive")
	notify := resp.(http.CloseNotifier).CloseNotify()
	logCh := make(chan string,255)
	go func() {
		for i:=1;i<1000;i++ {
			time.Sleep(1*time.Second) //not real code
			logCh <- version
		}
	}()
	flusher, ok := resp.(http.Flusher)
	if !ok {
		log.Println("Streaming not supported")
	}
	for {
		select {
		case <-notify:
			return
		case logs := <-logCh:
			fmt.Fprintln(resp, logs)
			flusher.Flush()
		}
	}
}

func (a Api) snap(c *cart.Context) {
	c.JSON(200, cart.H{
		"SnapShot": a.server.SnapShot,
	})
}

func (a Api) nodes(c *cart.Context) {
	c.JSON(200, cart.H{
		"Nodes": a.server.Nodes,
	})
}
