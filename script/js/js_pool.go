package js

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/script"
)

var (
	_ script.Client = &Pool{}
)

func NewPool(num int, option ...func(c script.Client)) *Pool {
	length := conv.Select[int](num <= 0, 1, num)
	length = conv.Select[int](length <= 0, 1, length)
	p := &Pool{
		length: length,
		list:   make([]*Client, 0, length),
		queue:  make(chan *Client, length),
	}
	for i := 0; i < p.length; i++ {
		c := New(option...)
		p.list = append(p.list, c)
		p.queue <- c
	}
	return p
}

type Pool struct {
	length int
	list   []*Client
	queue  chan *Client
	tag    *maps.Safe
}

func (this *Pool) Tag() *maps.Safe {
	if this.tag == nil {
		this.tag = maps.NewSafe()
	}
	return this.tag
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
