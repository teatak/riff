package riff

import (
	"encoding/json"
	"fmt"
	"github.com/gimke/cart"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpAPI struct{}

func (a *httpAPI) Index(r *cart.Router) {
	r.Route("/").GET(func(c *cart.Context) {
		c.Redirect(302, "/console/")
	})
	r.Route("/api", a.apiIndex)
}

func (a *httpAPI) apiIndex(r *cart.Router) {
	r.ANY(func(c *cart.Context, next cart.Next) {
		var q string
		if c.Request.Method == "GET" {
			q = c.Request.URL.Query().Get("query")
		}
		if c.Request.Method == "POST" {
			content := c.Request.Header.Get("content-type")
			if strings.ToLower(content) == "application/json" {
				b, _ := ioutil.ReadAll(c.Request.Body)
				var query map[string]interface{}
				json.Unmarshal(b, &query)
				q = query["query"].(string)
			} else if strings.ToLower(content) == "application/x-www-form-urlencoded" {
				q = c.Request.FormValue("query")
			}

		}
		result := executeQuery(q, schema)
		if len(result.Errors) > 0 {
			c.IndentedJSON(500, result)
		} else {
			c.IndentedJSON(200, result)
		}
	})
	//r.Route("").GET(a.version)
	//r.Route("/version").GET(a.version)
	//r.Route("/snap").GET(a.snap)
	//r.Route("/nodes").GET(a.nodes)
	//r.Route("/node/:name").GET(a.node)
	//r.Route("/services").GET(a.services)
	//r.Route("/service/:name", func(router *cart.Router) {
	//	router.Route("").GET(a.service)
	//	router.Route("/:command").GET(a.service)
	//	router.Route("/:command").POST(a.serviceCommand)
	//
	//})
	r.Route("/logs").GET(a.logs)
}

//func (a httpAPI) version(c *cart.Context) {
//	version := fmt.Sprintf("Cart version %s Riff version %s, build %s-%s", cart.Version, common.Version, common.GitBranch, common.GitSha)
//	c.IndentedJSON(200, cart.H{
//		"version": version,
//	})
//}
//
//func (a httpAPI) snap(c *cart.Context) {
//	c.IndentedJSON(200, cart.H{
//		"snapShot": server.SnapShot,
//	})
//}
//
//func (a httpAPI) nodes(c *cart.Context) {
//	c.IndentedJSON(200, server.api.Nodes())
//}
//
//func (a httpAPI) node(c *cart.Context) {
//	name, _ := c.Param("name")
//	c.IndentedJSON(200, server.api.Node(name))
//}
//
//func (a httpAPI) services(c *cart.Context) {
//	c.IndentedJSON(200, server.api.Services())
//}
//
//func (a httpAPI) service(c *cart.Context) {
//	name, _ := c.Param("name")
//	command, _ := c.Param("command")
//	state := api.StateAll
//	switch command {
//	case "alive":
//		state = api.StateAlive
//		break
//	case "dead":
//		state = api.StateDead
//		break
//	case "suspect":
//		state = api.StateSuspect
//		break
//	case "all":
//		state = api.StateAll
//		break
//	}
//	//!= "alive"
//	s := server.api.Service(name, state)
//	if s != nil {
//		c.IndentedJSON(200, s)
//	} else {
//		err := fmt.Errorf("not found")
//		c.IndentedJSON(404, cart.H{
//			"status": 404,
//			"error":  err.Error(),
//		})
//	}
//}
//
//func (a httpAPI) serviceCommand(c *cart.Context) {
//	name, _ := c.Param("name")
//	command, _ := c.Param("command")
//	var ok bool
//	switch command {
//	case "start":
//		ok = server.api.Start(name)
//		break
//	case "stop":
//		ok = server.api.Stop(name)
//		break
//	case "restart":
//		ok = server.api.Restart(name)
//		break
//	default:
//		err := fmt.Errorf("command missing")
//		c.IndentedJSON(400, cart.H{
//			"status": 400,
//			"error":  err.Error(),
//		})
//		return
//	}
//
//	switch ok {
//	case true:
//		c.IndentedJSON(200, cart.H{
//			"status": 200,
//		})
//		return
//	case false:
//		err := fmt.Errorf("not found")
//		c.IndentedJSON(404, cart.H{
//			"status": 404,
//			"error":  err.Error(),
//		})
//		return
//	}
//
//}

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
