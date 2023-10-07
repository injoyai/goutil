package js

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/script"
)

var (
	_ script.Interface = &Pool{}
)

func NewPool(num ...int) *Pool {
	length := conv.GetDefaultInt(1, num...)
	length = conv.SelectInt(length > 0, length, 1)
	return &Pool{
		length:  length,
		clients: make([]*Client, 0, length),
		queue:   make(chan *Client, length),
	}
}

type Pool struct {
	length  int
	clients []*Client
	queue   chan *Client
}

func (this *Pool) getClient() *Client {
	return <-this.queue
}

func (this *Pool) Exec(text string) (*conv.Var, error) {
	c := this.getClient()
	return c.Exec(text)
}

func (this *Pool) GetVar(key string) *conv.Var {
	c := this.getClient()
	return c.GetVar(key)
}

func (this *Pool) Set(key string, value interface{}) error {
	for _, v := range this.clients {
		if err := v.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (this *Pool) SetFunc(key string, value script.Func) error {
	return this.Set(key, value)
}

func (this *Pool) Close() error {
	for _, v := range this.clients {
		if err := v.Close(); err != nil {
			return err
		}
	}
	return nil
}
