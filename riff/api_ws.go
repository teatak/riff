package riff

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/teatak/cart"
	"math"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = math.MaxInt64
)

func (h *Http) handleWs(c *cart.Context) {
	ws, err := upgrader.Upgrade(c.Response, c.Request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			server.Logger.Printf(errorServicePrefix+"wrong ws connect: %v\n", err)
		}
		return
	}
	clientGone := make(chan bool)
	go h.handleWriter(ws, clientGone)
	go h.handleReader(ws, clientGone)
}

func (h *Http) handleWriter(ws *websocket.Conn, clientGone chan bool) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-ticker.C:
			h.mu.Lock()
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			h.mu.Unlock()
			if err != nil {
				server.Logger.Printf(infoServerPrefix+"wrong ws ping: %v\n", err)
				return
			}
		case <-clientGone:
			ticker.Stop()
			return
		}
	}
}

type WsResponse struct {
	Event string      `json:"event"`
	Body  interface{} `json:"body"`
}

type BodyWatch struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Query string `json:"query"`
}
type WsRequest struct {
	Event string    `json:"event"`
	Body  BodyWatch `json:"body"`
}

func (h *Http) handleReader(ws *websocket.Conn, clientGone chan bool) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	var watchHandler *httpWatchHandler
	var logHandler *httpLogHandler

	for {
		var request WsRequest
		err := ws.ReadJSON(&request)
		if err != nil {
			if closeError, ok := err.(*websocket.CloseError); ok {
				if closeError.Code == 1001 {
					server.Logger.Printf(infoServerPrefix+"client gone %v\n", err)
				} else {
					server.Logger.Printf(errorServerPrefix+"error %v\n", err)
				}
			} else {
				server.Logger.Printf(errorServerPrefix+"error %v\n", err)
			}
			clientGone <- true
			return
		}
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		switch request.Event {
		case "Watch":
			watch := request.Body
			//clear handler
			if watchHandler != nil {
				close(watchHandler.exitCh)
			}
			watchHandler = h.buildHttpWatchHandler(watch.Name, watch.Type)
			go h.handleWatch(ws, watchHandler, watch.Query)
		case "Logs":
			if logHandler != nil {
				close(logHandler.exitCh)
			}
			logHandler = h.buildHttpLogHandler()
			go h.handleLog(ws, logHandler)
		}
	}
}

func (h *Http) buildHttpWatchHandler(name string, watch string) *httpWatchHandler {
	watchType := NodeChanged
	switch watch {
	case "node":
		watchType = NodeChanged
		break
	case "service":
		watchType = ServiceChanged
		break
	}
	if watchType == NodeChanged && name == "" {
		name = server.Self.Name
	}
	handler := &httpWatchHandler{
		WatchParam: &WatchParam{
			Name:      name,
			WatchType: watchType,
		},
		watchCh: make(chan bool, 255),
		exitCh:  make(chan bool, 1),
	}
	return handler
}

func (h *Http) handleWatch(ws *websocket.Conn, handler *httpWatchHandler, query string) {
	server.watch.RegisterHandler(handler)
	defer server.watch.DeregisterHandler(handler)
	//params := graphql.Params{
	//	Schema:        schema,
	//	RequestString: query,
	//}
	opts := &RequestOptions{
		Query: query,
	}

	for {
		select {
		case <-handler.watchCh:
			//result := graphql.Do(params)
			result := h.Schema.Exec(context.Background(), opts.Query, opts.OperationName, opts.Variables)
			if len(result.Errors) > 0 {
				server.Logger.Printf(errorServicePrefix+"wrong result errors: %v\n", result.Errors)
			}
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if handler.WatchParam.WatchType == NodeChanged {
				h.mu.Lock()
				err := ws.WriteJSON(&WsResponse{
					Event: "NodeChange",
					Body:  result,
				})
				h.mu.Unlock()
				if err != nil {
					server.Logger.Printf(errorServicePrefix+"write node change errors: %v\n", err)
					return
				}
			} else {
				h.mu.Lock()
				err := ws.WriteJSON(&WsResponse{
					Event: "ServiceChange",
					Body:  result,
				})
				h.mu.Unlock()
				if err != nil {
					server.Logger.Printf(errorServicePrefix+"write service change errors: %v\n", err)
					return
				}
			}
		case <-handler.exitCh:
			return
		}
	}
}

func (h *Http) buildHttpLogHandler() *httpLogHandler {
	handler := &httpLogHandler{
		logCh:  make(chan string, 255),
		exitCh: make(chan bool, 1),
	}
	return handler
}

func (h *Http) handleLog(ws *websocket.Conn, handler *httpLogHandler) {
	server.logWriter.RegisterHandler(handler)
	defer server.logWriter.DeregisterHandler(handler)
	for {
		select {
		case logs := <-handler.logCh:
			h.mu.Lock()
			err := ws.WriteJSON(&WsResponse{
				Event: "Logs",
				Body:  logs,
			})
			h.mu.Unlock()
			if err != nil {
				server.Logger.Printf(errorServicePrefix+"write logs errors: %v\n", err)
				return
			}
		case <-handler.exitCh:
			return
		}
	}
}
