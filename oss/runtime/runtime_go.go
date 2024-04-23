package runtime

import (
	"context"
	"github.com/injoyai/base/maps"
	"sync/atomic"
	"time"
)

/*



 */

func Go(f func(ctx context.Context, args ...interface{}), args ...interface{}) *GoItem {
	return DefaultGoManage.Go(f, args...)
}

func Try(f func(ctx context.Context, args ...interface{}), args ...interface{}) *GoItem {
	return DefaultGoManage.Try(f, args...)
}

var (
	DefaultGoManage = NewGoManage()
)

func NewGoManage() *GoManage {
	return &GoManage{
		m:   maps.NewSafe(),
		key: &atomic.Int64{},
	}
}

type GoManage struct {
	m     *maps.Safe    //协程集合
	len   int           //协程数量
	key   *atomic.Int64 //key生成器
	limit int           //现在协程数量
}

func (this *GoManage) SetLimit(limit int) *GoManage {
	this.limit = limit
	return this
}

func (this *GoManage) Len() int {
	return this.len
}

func (this *GoManage) Close() error {
	this.Range(func(key int64, v *GoItem) bool {
		v.Close()
		return true
	})
	return nil
}

func (this *GoManage) Range(f func(key int64, v *GoItem) bool) {
	this.m.Range(func(key, value interface{}) bool {
		return f(key.(int64), value.(*GoItem))
	})
}

func (this *GoManage) Try(f func(ctx context.Context, args ...interface{}), args ...interface{}) *GoItem {
	if this.limit > 0 && this.Len() >= this.limit {
		return nil
	}
	return this.Go(f, args...)
}

func (this *GoManage) Go(f func(ctx context.Context, args ...interface{}), args ...interface{}) *GoItem {
	key := this.key.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	item := &GoItem{
		Key:  key,
		f:    f,
		args: args,
		final: func() {
			this.m.Del(key)
			this.len--
		},
		closed: make(chan struct{}),
		cancel: cancel,
	}
	this.len++
	this.m.Set(key, item)
	go item.run(ctx)
	return item
}

type GoItem struct {
	Key      int64                                          //key
	StarTime time.Time                                      //开始时间
	f        func(ctx context.Context, args ...interface{}) //协程的函数
	args     []interface{}                                  //参数
	final    func()                                         //内部结束执行
	closed   chan struct{}                                  //结束信号
	cancel   context.CancelFunc                             //关闭
}

// Close 停止协程,不一定生效,主要是关闭上下文,看协程内部是否实现
func (this *GoItem) Close() error {
	this.cancel()
	return nil
}

// Done 协程结束信号
func (this *GoItem) Done() <-chan struct{} {
	return this.closed
}

// Since 协程运行时间
func (this *GoItem) Since() time.Duration {
	return time.Since(this.StarTime)
}

func (this *GoItem) run(ctx context.Context) {
	this.StarTime = time.Now()
	this.f(ctx, this.args...)
	this.final()
}
