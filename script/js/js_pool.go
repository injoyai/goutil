package js

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/script"
)

var (
	_ script.Interface = &Pool{}
)

func NewPool(num ...int) *Pool {
	length := conv.GetDefaultInt(20, num...)
	length = conv.SelectInt(length <= 0, 1, length)
	p := &Pool{
		length: length,
		list:   make([]*Client, 0, length),
		queue:  make(chan *Client, length),
	}
	for i := 0; i < p.length; i++ {
		c := New()
		p.list = append(p.list, c)
		p.queue <- c
	}
	return p
}

type Pool struct {
	length int
	list   []*Client
	queue  chan *Client
}

func (this *Pool) get() *Client {
	return <-this.queue
}

func (this *Pool) put(c *Client) {
	this.queue <- c
}

func (this *Pool) Exec(text string, option ...func(i script.Client)) (interface{}, error) {
	c := this.get()
	val, err := c.Exec(text, option...)
	this.put(c)
	return val, err
}

func (this *Pool) Set(key string, value interface{}) error {
	for _, v := range this.list {
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
	for _, v := range this.list {
		if err := v.Close(); err != nil {
			return err
		}
	}
	return nil
}
