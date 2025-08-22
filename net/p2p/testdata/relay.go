package main

import (
	"encoding/json"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/mux"
	"github.com/injoyai/goutil/net/p2p"
	"github.com/injoyai/logs"
)

func init() {
	logs.SetShowColor(false)
}

func main() {
	cache := maps.NewGeneric[string, *P2PInfo]()

	s := mux.New(mux.WithPort(11111))
	s.ALL("/ws", func(r *mux.Request) {
		id := r.GetString("id")
		logs.Info("ID:", id)
		if cache.Exist(id) {
			in.Fail("该ID已经存在")
		}
		ws := r.Websocket()
		cache.Set(id, &P2PInfo{
			ID:        id,
			Websocket: ws,
		})
		defer func() {
			cache.Del(id)
			ws.Close()
		}()

		for {

			bs, err := ws.ReadMessage()
			if err != nil {
				return
			}

			m := p2p.Message{}
			err = json.Unmarshal(bs, &m)
			if err != nil {
				return
			}
			m.From = id

			bs, err = json.Marshal(m)
			if err != nil {
				return
			}

			ws2, ok := cache.Get(m.To)
			if !ok {
				ws.Write(p2p.Message{Type: p2p.Error, Data: "对方不在线"}.Bytes())
				continue
			}

			if _, err = ws2.Write(bs); err != nil {
				ws.Write(p2p.Message{Type: p2p.Error, Data: err.Error()}.Bytes())
				continue
			}

		}
	})
	s.Run()
}

type P2PInfo struct {
	ID string
	*mux.Websocket
}
