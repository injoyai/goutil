package mux

import (
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	WSText   = websocket.TextMessage
	WSBinary = websocket.BinaryMessage
	WSClose  = websocket.CloseMessage
	WSPing   = websocket.PingMessage
	WSPong   = websocket.PongMessage
)

type Websocket = websocket.Conn

func (this *Request) Websocket() (*Websocket, error) {
	up := websocket.Upgrader{
		ReadBufferSize:  1024 * 2,
		WriteBufferSize: 1024 * 2,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return up.Upgrade(this.Writer, this.Request, nil)
}
