package p2p

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/logs"
	"github.com/pion/webrtc/v3"
)

func NewListen(relay *client.Client) ios.ListenFunc {
	return func() (ios.Listener, error) {
		return Listen(relay)
	}
}

func Listen(relay *client.Client) (ios.Listener, error) {

	l := &listen{
		clients: maps.NewGeneric[string, *webrtc.PeerConnection](),
		accept:  chans.NewSafe[*Client](10),
		relay:   relay,
	}

	relay.OnDealMessage = func(c *client.Client, msg ios.Acker) {
		logs.PrintErr(func() error {
			m := Message{}
			err := json.Unmarshal(msg.Payload(), &m)
			if err != nil {
				return err
			}

			conn, err := l.clients.GetOrSetByHandler(m.From, func() (*webrtc.PeerConnection, error) {
				conn, err := webrtc.NewPeerConnection(webrtc.Configuration{})
				if err != nil {
					return nil, err
				}
				conn.OnDataChannel(func(channel *webrtc.DataChannel) {
					p := &Client{
						key:   m.From,
						ch:    chans.NewSafe[[]byte](100),
						dc:    channel,
						offer: webrtc.SessionDescription{},
					}
					channel.OnOpen(func() { l.accept.Must(p) })
					channel.OnMessage(func(msg webrtc.DataChannelMessage) {
						select {
						case p.ch.Chan <- msg.Data:
						default:
						}
					})
				})
				return conn, nil
			})
			if err != nil {
				return err
			}

			switch m.Type {
			case ICE:
				return conn.AddICECandidate(webrtc.ICECandidateInit{Candidate: m.Data})

			case SDP:
				desc := webrtc.SessionDescription{}
				err := json.Unmarshal([]byte(m.Data), &desc)
				if err != nil {
					return err
				}
				_, err = desc.Unmarshal()
				if err != nil {
					return err
				}
				if desc.Type == webrtc.SDPTypeOffer {
					err = conn.SetRemoteDescription(desc)
					if err != nil {
						return err
					}
					answer, err := conn.CreateAnswer(nil)
					if err != nil {
						return err
					}
					err = conn.SetLocalDescription(answer)
					if err != nil {
						return err
					}
					//监听候选地址,向中继服务器发送,由中继转发给目标
					conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
						if candidate == nil {
							return
						}
						relay.Write(Message{Type: ICE, To: m.From, Data: candidate.ToJSON().Candidate}.Bytes())
					})
					_, err = c.Write(Message{Type: SDP, To: m.From, Data: answer.SDP}.Bytes())
					return err
				}

			}
			return nil
		}())
	}
	go relay.Run(context.Background())

	return l, nil
}

var _ ios.Listener = &listen{}

type listen struct {
	clients *maps.Generic[string, *webrtc.PeerConnection]
	accept  *chans.Safe[*Client]
	relay   *client.Client
}

func (this *listen) Close() error {
	this.relay.Close()
	this.accept.Close()
	return nil
}

func (this *listen) Accept() (ios.ReadWriteCloser, string, error) {
	r, ok := <-this.accept.Chan
	if !ok {
		return nil, "", errors.New("listener closed")
	}
	return r, r.key, nil
}

func (this *listen) Addr() string {
	return fmt.Sprintf("%p", this)
}
