package riff

import (
	"encoding/json"
	"fmt"
	"github.com/gimke/cart"
	"github.com/graphql-go/graphql"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type httpAPI struct{}

func (a *httpAPI) Index(r *cart.Router) {
	r.Route("/").GET(func(c *cart.Context) {
		c.Redirect(302, "/console/")
	})
	r.Route("/api", a.apiIndex)
}

type RequestOptions struct {
	Query         string                 `json:"query" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName" url:"operationName" schema:"operationName"`
}

// a workaround for getting`variables` as a JSON string
type requestOptionsCompatibility struct {
	Query         string `json:"query" url:"query" schema:"query"`
	Variables     string `json:"variables" url:"variables" schema:"variables"`
	OperationName string `json:"operationName" url:"operationName" schema:"operationName"`
}

const (
	ContentTypeJSON           = "application/json"
	ContentTypeGraphQL        = "application/graphql"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

func getFromForm(values url.Values) *RequestOptions {
	query := values.Get("query")
	if query != "" {
		// get variables map
		variables := make(map[string]interface{}, len(values))
		variablesStr := values.Get("variables")
		json.Unmarshal([]byte(variablesStr), &variables)

		return &RequestOptions{
			Query:         query,
			Variables:     variables,
			OperationName: values.Get("operationName"),
		}
	}

	return nil
}

func NewRequestOptions(r *http.Request) *RequestOptions {
	if reqOpt := getFromForm(r.URL.Query()); reqOpt != nil {
		return reqOpt
	}

	if r.Method != "POST" {
		return &RequestOptions{}
	}

	if r.Body == nil {
		return &RequestOptions{}
	}

	// TODO: improve Content-Type handling
	contentTypeStr := r.Header.Get("Content-Type")
	contentTypeTokens := strings.Split(contentTypeStr, ";")
	contentType := contentTypeTokens[0]

	switch contentType {
	case ContentTypeGraphQL:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &RequestOptions{}
		}
		return &RequestOptions{
			Query: string(body),
		}
	case ContentTypeFormURLEncoded:
		if err := r.ParseForm(); err != nil {
			return &RequestOptions{}
		}

		if reqOpt := getFromForm(r.PostForm); reqOpt != nil {
			return reqOpt
		}

		return &RequestOptions{}

	case ContentTypeJSON:
		fallthrough
	default:
		var opts RequestOptions
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &opts
		}
		err = json.Unmarshal(body, &opts)
		if err != nil {
			// Probably `variables` was sent as a string instead of an object.
			// So, we try to be polite and try to parse that as a JSON string
			var optsCompatible requestOptionsCompatibility
			json.Unmarshal(body, &optsCompatible)
			json.Unmarshal([]byte(optsCompatible.Variables), &opts.Variables)
		}
		return &opts
	}
}

func (a *httpAPI) apiIndex(r *cart.Router) {
	r.ANY(func(c *cart.Context, next cart.Next) {
		//var reqOpt *RequestOptions
		opts := NewRequestOptions(c.Request)

		params := graphql.Params{
			Schema:         schema,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			Context:        c.Request.Context(),
		}
		result := executeQuery(params)
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
