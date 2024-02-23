package trunk

import "github.com/injoyai/base/maps"

func New() *Entity {
	return &Entity{
		Safe: maps.NewSafe(),
	}
}

type Entity struct {
	*maps.Safe
}

func (this *Entity) get(topic string) *Mem {
	if this.Safe == nil {
		this.Safe = maps.NewSafe()
	}
	e, _ := this.GetOrSetByHandler(topic, func() (interface{}, error) {
		m := NewMem()
		go m.Run()
		return m, nil
	})
	return e.(*Mem)
}

func (this *Entity) Publish(topic string, data interface{}) {
	this.get(topic).Publish(data)
}

func (this *Entity) Subscribe(topic string, h SubscribeHandler) {
	this.get(topic).Subscribe(h)
}

func (this *Entity) Middleware(topic string, h MiddlewareHandler) {
	this.get(topic).Middleware(h)
}

func (this *Entity) Hook(topic string) *Hook {
	return this.get(topic).Hook()
}
