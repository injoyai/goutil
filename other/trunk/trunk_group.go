package trunk

import "github.com/injoyai/base/maps"

func NewGroup() *Group {
	return &Group{
		Safe: maps.NewSafe(),
	}
}

type Group struct {
	*maps.Safe
}

func (this *Group) Publish(topic string, data interface{}) {
	this.get(topic).Publish(data)
}

func (this *Group) Subscribe(topic string, h func(msg interface{})) {
	this.get(topic).Subscribe(h)
}

func (this *Group) Middleware(topic string, h func(msg interface{}) bool) {
	this.get(topic).Middleware(h)
}

func (this *Group) Hook(topic string, cap uint) *Channel {
	return this.get(topic).Channel(cap)
}

func (this *Group) get(topic string) *Trunk {
	if this.Safe == nil {
		this.Safe = maps.NewSafe()
	}
	e, _ := this.GetOrSetByHandler(topic, func() (interface{}, error) {
		return New(), nil
	})
	return e.(*Trunk)
}
