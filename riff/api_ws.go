package riff

import (
	"github.com/gimke/cart"
	"github.com/gorilla/websocket"
	"time"
	"github.com/graphql-go/graphql"
	"fmt"
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
			fmt.Println("ping")
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}

func (h *Http) handleReader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	name := ""
	watch := ""
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

	params := graphql.Params{
		Schema:         schema,
		RequestString:  `{
  node(name:"node1") {
    name
  }
}`,
	}
	go func() {
		for{
			select {
			case <-handler.serviceCh:
				result := graphql.Do(params)
				if len(result.Errors) > 0 {
					server.Logger.Printf(errorServicePrefix+"wrong result, unexpected errors: %v\n", result.Errors)
				}
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				ws.WriteJSON(result)
				//b, _ := json.Marshal(result)
			}
		}
	}()


	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		ws.WriteJSON(struct {
			Name string
		}{
			Name:"aaa",
		})
	}
}