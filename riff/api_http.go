package riff

import (
	"encoding/json"
	"fmt"
	"github.com/gimke/cart"
	//graphql1 "github.com/graphql-go/graphql"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"github.com/graph-gophers/graphql-go"
)

const (
	ContentTypeJSON           = "application/json"
	ContentTypeGraphQL        = "application/graphql"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

type Http struct {
	mu sync.Mutex
	Schema *graphql.Schema
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

func (h *Http) Index(r *cart.Router) {
	r.Route("/").GET(func(c *cart.Context) {
		c.Redirect(302, "/console/")
	})
	r.Route("/api", h.apiIndex)
	r.Route("/ws").GET(h.handleWs)
}

func (h *Http) apiIndex(r *cart.Router) {
	r.ANY(h.api)
	r.Route("/watch").ANY(h.watch)
	r.Route("/logs").ANY(h.logs)
}

func (h *Http) getFromForm(values url.Values) *RequestOptions {
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

func (h *Http) newRequestOptions(r *http.Request) *RequestOptions {
	if reqOpt := h.getFromForm(r.URL.Query()); reqOpt != nil {
		return reqOpt
	}

	if r.Method != "POST" {
		return &RequestOptions{}
	}

	if r.Body == nil {
		return &RequestOptions{}
	}

	contentTypeStr := r.Header.Get("Content-Type")
	contentTypeTokens := strings.Split(contentTypeStr, ";")
	contentType := strings.TrimSpace(contentTypeTokens[0])

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

		if reqOpt := h.getFromForm(r.PostForm); reqOpt != nil {
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

func (h *Http) api(c *cart.Context, next cart.Next) {
	//var reqOpt *RequestOptions
	opts := h.newRequestOptions(c.Request)
	//params := graphql1.Params{
	//	Schema:         schema,
	//	RequestString:  opts.Query,
	//	VariableValues: opts.Variables,
	//	OperationName:  opts.OperationName,
	//	Context:        c.Request.Context(),
	//}
	//result := graphql1.Do(params)
	result := h.Schema.Exec(c.Request.Context(), opts.Query, opts.OperationName, opts.Variables)

	if len(result.Errors) > 0 {
		server.Logger.Printf(errorServicePrefix+"wrong result, unexpected errors: %v\n", result.Errors)
		c.IndentedJSON(500, result)
	} else {
		c.IndentedJSON(200, result)
	}
}

type httpLogHandler struct {
	logCh  chan string
	exitCh chan bool
}

func (h *httpLogHandler) HandleLog(log string) {
	// Do a non-blocking send
	select {
	case h.logCh <- log:
	}
}

func (h *Http) logs(c *cart.Context, next cart.Next) {
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
