package runtimes

import (
	"time"
)

// GoPool 协程复用
type GoPool struct {
	ch          chan func()
	pool        chan struct{}
	IdleTimeout time.Duration //空闲协程超时时间,超时会退出协程
}

func NewGoPool(limit uint) *GoPool {
	if limit <= 0 {
		limit = 1
	}
	p := &GoPool{
		ch:          make(chan func()),
		pool:        make(chan struct{}, limit),
		IdleTimeout: time.Second * 10,
	}
	return p
}

func (this *GoPool) Current() int {
	return len(this.pool)
}

func (this *GoPool) Go(f func()) {
	select {
	case this.ch <- f:
		//尝试加入执行队列,当有协程空闲时,会被立马消费,否则会阻塞
		//当队列阻塞时,尝试申请新协程,新协程执行会立马消费一个任务
	case this.pool <- struct{}{}:
		go this.run()
	}
}

func (this *GoPool) run() error {
	t := time.NewTimer(this.IdleTimeout)
	defer func() {
		t.Stop()
		<-this.pool
	}()
	for {
		select {
		case <-t.C:
			return nil
		case f, ok := <-this.ch:
			if !ok {
				return nil
			}
			if f != nil {
				f()
			}
		}
	}
}
