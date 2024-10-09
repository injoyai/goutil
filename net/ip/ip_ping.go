package ip

import (
	"context"
	"net"
	"time"
)

func Ping(ip string, timeout time.Duration) (time.Duration, error) {
	conn, err := net.DialTimeout("ip:icmp", ip, timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	t := time.Now()
	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return 0, err
	}
	if _, err = conn.Write([]byte{
		8, 0, 247, 253, 0, 1, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
		return 0, err
	}
	buf := make([]byte, 65535)
	_, err = conn.Read(buf)
	return time.Since(t), err
}

func NewPinger() *Pinger {
	return &Pinger{
		Host: []string{
			"114.114.114.114", //运营商 约5ms
			"www.baidu.com",   //百度 约30ms
			"www.aliyun.com",  //阿里 约5ms
			"www.tencent.com", //腾讯 约10ms
		},
		i:   0,
		c:   nil,
		t:   time.Second,
		buf: make([]byte, 1024),
	}
}

type Pinger struct {
	Host []string
	i    int
	c    net.Conn
	t    time.Duration
	buf  []byte
}

func (this *Pinger) SetTimeout(t time.Duration) {
	this.t = t
}

func (this *Pinger) SetHost(host []string) {
	this.Host = host
}

func (this *Pinger) For(ctx context.Context, interval time.Duration, f func(time.Duration, error)) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			s, err := this.Ping()
			f(s, err)
		}
	}
}

func (this *Pinger) Ping() (time.Duration, error) {
	if this.c == nil {
		conn, err := net.DialTimeout("ip:icmp", this.Host[this.i%len(this.Host)], this.t)
		if err != nil {
			return 0, err
		}
		this.c = conn
		this.i++
	}
	t := time.Now()
	if err := this.c.SetDeadline(time.Now().Add(this.t)); err != nil {
		this.c = nil
		return 0, err
	}
	if _, err := this.c.Write([]byte{
		8, 0, 247, 253, 0, 1, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
		return 0, err
	}
	if _, err := this.c.Read(this.buf); err != nil {
		this.c = nil
		return 0, err
	}
	return time.Since(t), nil
}
