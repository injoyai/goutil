package terminal

import (
	"encoding/base64"

	ws "github.com/gorilla/websocket"
	"github.com/injoyai/conv"
	"github.com/injoyai/ios/v2"
	"github.com/injoyai/ios/v2/client"
	"github.com/injoyai/ios/v2/client/dial"
	"github.com/injoyai/ios/v2/module/ssh"
	json "github.com/json-iterator/go"
)

type Websocket interface {
	// Write 写入数据,写入\n执行
	Write(p []byte) (int, error)
	// ReadMessage 读取拆包的数据
	ReadMessage() ([]byte, error)
}

func NewWebsocket(s *ws.Conn, cfg *ssh.Config, options ...client.Option) (*websocket, error) {
	c, err := dial.With(ssh.NewDial(cfg), options...)
	if err != nil {
		return nil, err
	}
	c.OnDealMessage(func(c *client.Client, msg ios.Acker) {
		if err := s.WriteMessage(ws.TextMessage, conv.Bytes(WebsocketMsg{
			Type: WsTypeCmd,
			Data: base64.StdEncoding.EncodeToString(msg.Bytes()),
		})); err != nil {
			c.Close()
		}
	})
	c.OnDisconnect(func(c *client.Client, err error) {
		s.WriteMessage(ws.TextMessage, conv.Bytes(&WebsocketMsg{
			Type: WsTypeErr,
			Data: err.Error(),
		}))
		s.Close()
	})
	return &websocket{
		Client: c,
		ws:     s,
	}, nil
}

type websocket struct {
	*client.Client
	ws *ws.Conn
}

func (this *websocket) Run() error {
	go this.Client.Run()
	for {
		_, data, err := this.ws.ReadMessage()
		if err != nil {
			return err
		}

		msg := new(WebsocketMsg)
		json.Unmarshal(data, msg)

		switch msg.Type {

		case WsTypeResize:
			//重新设置窗口大小
			if msg.High > 0 && msg.Wide > 0 {
				if err := this.Client.Origin().(*ssh.Client).WindowChange(msg.High, msg.Wide); err == nil {
					if err := this.ws.WriteMessage(ws.TextMessage, data); err != nil {
						return err
					}
				}
			}

		case WsTypeHeartbeat:
			//心跳数据,原路返回
			if err := this.ws.WriteMessage(ws.TextMessage, data); err != nil {
				return err
			}

		case WsTypeCmd:
			//命令
			decodeBytes, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				if err := this.ws.WriteMessage(ws.TextMessage, conv.Bytes(WebsocketMsg{
					Type: WsTypeErr,
					Data: err.Error(),
				})); err != nil {
					return err
				}
			}
			if _, err := this.Client.Write(decodeBytes); err == nil {
				if err := this.ws.WriteMessage(ws.TextMessage, conv.Bytes(WebsocketMsg{
					Type: WsTypeErr,
					Data: err.Error(),
				})); err != nil {
					return err
				}
			}
		}

	}
}

const (
	WsTypeCmd       = "cmd"
	WsTypeResize    = "resize"
	WsTypeHeartbeat = "heartbeat"
	WsTypeErr       = "err"
)

type WebsocketMsg struct {
	Type      string `json:"type"`
	Data      string `json:"data,omitempty"`      // WsTypeCmd,WsTypeErr
	High      int    `json:"high,omitempty"`      // WsTypeResize
	Wide      int    `json:"wide,omitempty"`      // WsTypeResize
	Timestamp int    `json:"timestamp,omitempty"` // WsTypeHeartbeat
}
