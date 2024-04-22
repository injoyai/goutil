package runtime

import (
	"context"
	"github.com/injoyai/base/maps"
	"sync/atomic"
	"time"
)

/*



 */

func Go(f func(ctx context.Context)) {
	DefaultGoManage.Go(f)
}

func GoLen() int64 {
	return DefaultGoManage.len
}

var (
	DefaultGoManage = &_goManage{m: maps.NewSafe()}
)

type _goManage struct {
	m   *maps.Safe   //协程集合
	len int64        //协程数量
	key atomic.Int64 //
}

func (this *_goManage) Go(f func(ctx context.Context)) {
	key := this.key.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	v := &_go{
		key: key,
		f:   f,
		final: func() {
			this.m.Del(key)
			this.len--
		},
		closed: make(chan struct{}),
		cancel: cancel,
	}
	this.len++
	this.m.Set(key, v)
	go v.run(ctx)
}

type _go struct {
	key      int64     //key
	StarTime time.Time //开始时间
	f        func(ctx context.Context)
	final    func()
	closed   chan struct{}
	cancel   context.CancelFunc
}

func (this *_go) Close() error {
	this.cancel()
	return nil
}

func (this *_go) Done() <-chan struct{} {
	return this.closed
}

func (this *_go) Since() time.Duration {
	return time.Since(this.StarTime)
}

func (this *_go) run(ctx context.Context) {
	this.StarTime = time.Now()
	this.f(ctx)
	this.final()
}
