package p2p

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/pion/webrtc/v3"
	"io"
	"time"
)

func NewDial(relay *client.Client, target string) ios.DialFunc {
	return func(ctx context.Context) (ios.ReadWriteCloser, string, error) {
		c, err := Dial(relay, target)
		return c, target, err
	}
}

func Dial(relay *client.Client, target string) (*Client, error) {
	conn, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return nil, err
	}

	//创建数据通道,v1版本先默认创建一个数据通道,先执行这个才能收集ICE信息
	dc, err := conn.CreateDataChannel("data", nil)
	if err != nil {
		return nil, err
	}

	offer, err := conn.CreateOffer(nil)
	if err != nil {
		return nil, err
	}

	err = conn.SetLocalDescription(offer)
	if err != nil {
		return nil, err
	}

	//// ✅ 用这个通道注册等待 ICE 收集完成
	<-webrtc.GatheringCompletePromise(conn)

	//像中继服务器发送请求连接数据
	_, err = relay.Write(Message{Type: SDP, To: target, Data: conv.String(conn.LocalDescription())}.Bytes())
	if err != nil {
		return nil, err
	}
	wait := chans.NewSafe[struct{}]()
	dc.OnOpen(func() { wait.Add(struct{}{}) })
	relay.OnDealMessage = func(c *client.Client, msg ios.Acker) {

		var err error
		defer func() { wait.CloseWithErr(err) }()

		m := Message{}
		err = json.Unmarshal(msg.Payload(), &m)
		if err != nil {
			return
		}

		switch m.Type {
		case SDP:
			desc := webrtc.SessionDescription{}
			err = json.Unmarshal([]byte(m.Data), &desc)
			if err != nil {
				return
			}
			_, err = desc.Unmarshal()
			if err != nil {
				return
			}
			//设置远程配置描述
			if desc.Type == webrtc.SDPTypeAnswer {
				err = conn.SetRemoteDescription(desc)
				return
			}
		case Error:
			err = errors.New(m.Data)
		}

	}

	//等待中继服务器响应
	select {
	case <-wait.Done():
		return nil, wait.Err()

	case <-time.After(time.Second * 10):
		return nil, errors.New("建立连接超时")

	case <-wait.Chan:

	}

	p := &Client{
		key:   target,
		ch:    chans.NewSafe[[]byte](100),
		dc:    dc,
		offer: offer,
	}

	//监听候选地址,向中继服务器发送,由中继转发给目标
	conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}
		relay.Write(Message{Type: ICE, To: target, Data: candidate.ToJSON().Candidate}.Bytes())
	})

	p.dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		select {
		case p.ch.Chan <- msg.Data:
		default:
		}
	})

	return p, nil
}

var _ ios.MReadWriteCloser = &Client{}

type Client struct {
	key   string
	ch    *chans.Safe[[]byte]
	dc    *webrtc.DataChannel
	offer webrtc.SessionDescription
}

func (this *Client) ReadMessage() ([]byte, error) {
	if this.ch.Closed() {
		return nil, io.EOF
	}
	bs, ok := <-this.ch.Chan
	if !ok {
		return nil, io.EOF
	}
	return bs, nil
}

func (this *Client) Write(p []byte) (int, error) {
	err := this.dc.Send(p)
	return len(p), err
}

func (this *Client) WriteString(s string) (int, error) {
	err := this.dc.SendText(s)
	return len(s), err
}

func (this *Client) Close() error {
	this.ch.Close()
	return this.dc.Close()
}
