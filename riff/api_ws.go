package riff

import (
	"github.com/gimke/cart"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (h *Http) handleWs(c *cart.Context) {
	ws, err := upgrader.Upgrade(c.Response, c.Request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			server.Logger.Printf(errorServicePrefix+"wrong ws, unexpected errors: %v\n", err)
		}
		return
	}
	go h.handleWriter(ws)
	h.handleReader(ws)
}

func (h *Http) handleWriter(ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
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

func (h *Http) handleReader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	//name := server.Self.Name
	//handler := h.buildHttpServiceHandler(name,"node")

	var handler *httpServiceHandler

	for {
		var request WsRequest
		err := ws.ReadJSON(&request)
		if err != nil {
			server.Logger.Printf(errorServerPrefix+"error %v\n", err)
			break
		}
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		switch request.Event {
		case "Watch":
			watch := request.Body
			//clear handler
			if handler != nil {
				close(handler.exitCh)
			}
			handler = h.buildHttpServiceHandler(watch.Name, watch.Type)
			go h.handleWatch(ws, handler, watch.Query)
			//ws.WriteJSON()
		}
	}
}

func (h *Http) buildHttpServiceHandler(name string, watch string) *httpServiceHandler {

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

	handler := &httpServiceHandler{
		WatchParam: &WatchParam{
			Name:      name,
			WatchType: watchType,
		},
		serviceCh: make(chan bool, 512),
		exitCh:    make(chan bool, 1),
	}
	return handler
}

func (h *Http) handleWatch(ws *websocket.Conn, handler *httpServiceHandler, query string) {

	server.watch.RegisterHandler(handler)
	defer server.watch.DeregisterHandler(handler)

	params := graphql.Params{
		Schema:        schema,
		RequestString: query,
	}
	for {
		select {
		case <-handler.serviceCh:
			result := graphql.Do(params)
			if len(result.Errors) > 0 {
				server.Logger.Printf(errorServicePrefix+"wrong result, unexpected errors: %v\n", result.Errors)
			}
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if handler.WatchParam.WatchType == NodeChanged {
				h.mu.Lock()
				ws.WriteJSON(&WsResponse{
					Event: "NodeChange",
					Body:  result,
				})
				h.mu.Unlock()
			} else {
				h.mu.Lock()
				ws.WriteJSON(&WsResponse{
					Event: "ServiceChange",
					Body:  result,
				})
				h.mu.Unlock()
			}

		case <-handler.exitCh:
			server.Logger.Printf(infoServicePrefix + "exit handler\n")
			return
			//b, _ := json.Marshal(result)
		}
	}
}
