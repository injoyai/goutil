package main

import (
	"context"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/goutil/net/p2p"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/client/redial"
	"github.com/injoyai/logs"
	"time"
)

func init() {
	cfg.WithFlag(
		&cfg.Flag{Name: "id", Default: "A"},
		&cfg.Flag{Name: "to", Default: "B"},
	)
}

func main() {

	id := cfg.GetString("id", "A")
	to := cfg.GetString("to", "Server")

	ws := redial.Websocket("ws://39.107.120.124:11111/ws?id="+id, func(c *client.Client) {
		//c.Logger.Debug(false)
		//c.Logger.SetLevel(common.LevelError)
		c.OnWrite = client.NewWriteSafe()
	})
	go ws.Run(context.Background())

	err := redial.Run(
		p2p.NewDial(ws, to),

		func(c *client.Client) {
			c.OnReconnect = client.NewReconnectInterval(time.Second * 8)
			c.OnConnected = func(c *client.Client) error {
				c.GoTimerWriter(time.Second*5, func(w ios.MoreWriter) error {
					_, err := w.WriteString(time.Now().String())
					return err
				})
				return nil
			}
		},
	)

	logs.Err(err)

}
