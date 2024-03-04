package trunk

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func New() *Trunk {
	t := &Trunk{
		ch: make(chan interface{}, 1000),
	}
	go t.run()
	return t
}

type Trunk struct {
	ch          chan interface{}
	middlewares []*Middleware
	subscribes  []*Subscribe
	channels    []*Channel
	mu          sync.Mutex
}

func (this *Trunk) run() {
	wg := sync.WaitGroup{}
	for msg := range this.ch {
		for _, v := range this.channels {
			//尝试加入通道,加入失败则丢弃
			v.Publish(msg)
		}
		for _, v := range this.subscribes {
			wg.Add(1)
			go func(v *Subscribe) {
				defer func() {
					recover()
					wg.Done()
				}()
				v.f(msg)
			}(v)
		}
		wg.Wait()
	}
}

// Publish 发布数据到队列,并运行中间件
func (this *Trunk) Publish(msg interface{}) {
	for _, v := range this.middlewares {
		if !v.f(msg) {
			return
		}
	}
	this.ch <- msg
}

// Middleware 从队列中拦截数据,可以进行修改/过滤操作
func (this *Trunk) Middleware(f func(msg interface{}) bool) *Middleware {
	return &Middleware{
		t: this,
		k: fmt.Sprintf("%p", f),
		f: f,
	}
}

// Subscribe 订阅数据,可以同时被多个订阅
func (this *Trunk) Subscribe(f func(msg interface{})) *Subscribe {
	s := &Subscribe{
		t: this,
		k: fmt.Sprintf("%p", f),
		f: f,
	}
	this.subscribes = append(this.subscribes, s)
	return s
}

// Channel 接入一个通道到队列,是订阅的另一种方式
func (this *Trunk) Channel(cap uint) *Channel {
	c := make(chan interface{}, cap)
	return &Channel{
		t: this,
		k: fmt.Sprintf("%p", c),
		C: c,
	}
}

type Subscribe struct {
	t *Trunk
	k string
	f func(msg interface{})
}

func (this *Subscribe) Close() error {
	this.t.mu.Lock()
	defer this.t.mu.Unlock()
	for i, v := range this.t.subscribes {
		if v.k == this.k {
			this.t.subscribes = append(this.t.subscribes[:i], this.t.subscribes[i+1:]...)
			break
		}
	}
	return nil
}

type Middleware struct {
	t *Trunk
	k string
	f func(msg interface{}) bool
}

func (this *Middleware) Close() error {
	this.t.mu.Lock()
	defer this.t.mu.Unlock()
	for i, v := range this.t.middlewares {
		if v.k == this.k {
			this.t.middlewares = append(this.t.middlewares[:i], this.t.middlewares[i+1:]...)
			break
		}
	}
	return nil
}

type Channel struct {
	t      *Trunk
	k      string
	C      chan interface{}
	closed uint32
}

func (this *Channel) Publish(msg interface{}) {
	if atomic.LoadUint32(&this.closed) == 1 {
		//通道状态为1(关闭)，直接返回
		return
	}
	select {
	//尝试加入通道,加入失败则丢弃
	case this.C <- msg:
	default:
	}
}

func (this *Channel) Close() error {
	this.t.mu.Lock()
	defer this.t.mu.Unlock()
	for i, v := range this.t.channels {
		if v.k == this.k {
			this.t.channels = append(this.t.channels[:i], this.t.channels[i+1:]...)
			//设置通道状态为1(关闭)
			atomic.SwapUint32(&this.closed, 1)
			close(this.C)
			break
		}
	}
	return nil
}
