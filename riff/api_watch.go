package riff

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teatak/cart"
)

type httpWatchHandler struct {
	*WatchParam
	watchCh chan bool
	exitCh  chan bool
}

func (h *httpWatchHandler) HandleWatch() {
	// Do a non-blocking send
	h.watchCh <- true
}

func (h *httpWatchHandler) GetParam() *WatchParam {
	return h.WatchParam
}

func (h *Http) watch(c *cart.Context, next cart.Next) {
	resp := c.Response
	ctx := c.Request.Context()

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
	case "service":
		watchType = ServiceChanged
	}

	watchHandler := &httpWatchHandler{
		WatchParam: &WatchParam{
			Name:      name,
			WatchType: watchType,
		},
		watchCh: make(chan bool, 512),
	}
	server.watch.RegisterHandler(watchHandler)
	defer server.watch.DeregisterHandler(watchHandler)

	opts := h.newRequestOptions(c.Request)
	//params := graphql.Params{
	//	Schema:         schema,
	//	RequestString:  opts.Query,
	//	VariableValues: opts.Variables,
	//	OperationName:  opts.OperationName,
	//	Context:        c.Request.Context(),
	//}

	flusher, ok := resp.(http.Flusher)
	if !ok {
		server.Logger.Println("streaming not supported")
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-watchHandler.watchCh:

			//result := graphql.Do(params)
			result := h.Schema.Exec(c.Request.Context(), opts.Query, opts.OperationName, opts.Variables)
			if len(result.Errors) > 0 {
				server.Logger.Printf(errorServicePrefix+"wrong result, unexpected errors: %v\n", result.Errors)
			}
			b, _ := json.Marshal(result)
			fmt.Fprintln(resp, string(b))

			flusher.Flush()
		}
	}
}
