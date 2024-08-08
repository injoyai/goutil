package mux

import "github.com/gorilla/websocket"

const (
	WSText   = websocket.TextMessage
	WSBinary = websocket.BinaryMessage
	WSClose  = websocket.CloseMessage
	WSPing   = websocket.PingMessage
	WSPong   = websocket.PongMessage
)

type Websocket = websocket.Conn
