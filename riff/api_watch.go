package riff

import (
	"encoding/json"
	"fmt"
	"github.com/gimke/cart"
	"github.com/graphql-go/graphql"
	"net/http"
)

type httpServiceHandler struct {
	*WatchParam
	serviceCh chan bool
}

func (h *httpServiceHandler) HandleWatch() {
	// Do a non-blocking send
	select {
	case h.serviceCh <- true:
	}
}

func (h *httpServiceHandler) GetParam() *WatchParam {
	return h.WatchParam
}

func (h *Http) watch(c *cart.Context, next cart.Next) {
	resp := c.Response
	clientGone := resp.(http.CloseNotifier).CloseNotify()

	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.Header().Set("Connection", "Keep-Alive")
	resp.Header().Set("Transfer-Encoding", "chunked")
	resp.Header().Set("X-Content-Type-Options", "nosniff")

	//get type and name
	name := c.Request.URL.Query().Get("name")
	watch := c.Request.URL.Query().Get("type")
	watchType := NodeChanged
	if name == "" {
		name = server.Self.Name
	}
	switch watch {
	case "node":
		watchType = NodeChanged
		break
	case "service":
		watchType = ServiceChanged
		break
	}

	handler := &httpServiceHandler{
		WatchParam: &WatchParam{
			Name:      name,
			WatchType: watchType,
		},
		serviceCh: make(chan bool, 512),
	}
	server.watch.RegisterHandler(handler)
	defer server.watch.DeregisterHandler(handler)

	opts := h.newRequestOptions(c.Request)
	params := graphql.Params{
		Schema:         schema,
		RequestString:  opts.Query,
		VariableValues: opts.Variables,
		OperationName:  opts.OperationName,
		Context:        c.Request.Context(),
	}

	flusher, ok := resp.(http.Flusher)
	if !ok {
		server.Logger.Println("Streaming not supported")
	}
	for {
		select {
		case <-clientGone:
			return
		case <-handler.serviceCh:

			result := graphql.Do(params)
			if len(result.Errors) > 0 {
				server.Logger.Printf(errorServicePrefix+"wrong result, unexpected errors: %v\n", result.Errors)
			}
			b, _ := json.Marshal(result)
			fmt.Fprintln(resp, string(b))

			flusher.Flush()
		}
	}
}
