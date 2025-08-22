package main

import (
	"context"
	"github.com/injoyai/goutil/net/p2p"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/client/redial"
	"github.com/injoyai/ios/server"
)

func main() {
	id := "Server"

	ws := redial.Websocket("ws://39.107.120.124:11111/ws?id="+id, func(c *client.Client) {
		//c.Logger.Debug(false)
		c.OnWrite = client.NewWriteSafe()
	})
	go ws.Run(context.Background())

	server.Run(p2p.NewListen(id, ws))
}
