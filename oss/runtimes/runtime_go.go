package runtimes

import (
	"context"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/base/safe"
	"reflect"
	"runtime"
	"sync/atomic"
	"time"
)

type GoHandler func(ctx context.Context, args ...any)

func Go(f GoHandler, args ...any) *GoItem {
	return DefaultGoManage.Go(f, args...)
}

var (
	DefaultGoManage = NewGoManage()
)

func NewGoManage() *GoManage {
	return &GoManage{
		m: maps.NewSafe(),
	}
}

type GoManage struct {
	m       *maps.Safe //协程集合
	total   int64      //已运行的协程数量,作为协程的编号
	current int64      //正在运行的协程数量
	limit   int        //限制的协程数量
}

func (this *GoManage) SetLimit(limit int) *GoManage {
	this.limit = limit
	return this
}

func (this *GoManage) Len() int {
	return int(this.current)
}

func (this *GoManage) Close() error {
	this.Range(func(key int64, v *GoItem) bool {
		v.Stop()
		return true
	})
	return nil
}

func (this *GoManage) Range(f func(key int64, v *GoItem) bool) {
	this.m.Range(func(key, value any) bool {
		return f(key.(int64), value.(*GoItem))
	})
}

func (this *GoManage) Go(f GoHandler, args ...any) *GoItem {
	if this.limit > 0 && this.Len() >= this.limit {
		return nil
	}
	key := atomic.AddInt64(&this.total, 1)
	item := &GoItem{
		Key:    key,
		f:      f,
		Input:  args,
		Runner: safe.NewRunner(nil),
	}
	item.Runner.SetFunc(func(ctx context.Context) error {
		atomic.AddInt64(&this.current, 1)
		this.m.Set(key, item)
		item.Panic = item.run(ctx)
		this.m.Del(key)
		atomic.AddInt64(&this.current, -1)
		return item.Panic
	})
	item.Start()
	return item
}

type GoItem struct {
	Key         int64         //key
	StarTime    time.Time     //开始时间
	StopTime    time.Time     //结束信号时间,手动结束才会赋值
	StoppedTime time.Time     //结束完成时间
	f           GoHandler     //函数
	fn          *runtime.Func //函数信息
	Input       []any         //函数的参数
	Panic       error         //panic信息,正常未nil
	*safe.Runner
}

// todo 是否能获取到内存信息?

// Since 协程已经运行时间
func (this *GoItem) Since() time.Duration {
	return time.Since(this.StarTime)
}

func (this *GoItem) Func() *runtime.Func {
	if this.fn == nil {
		pc := reflect.ValueOf(this.f).Pointer()
		this.fn = runtime.FuncForPC(pc)
	}
	return this.fn
}

func (this *GoItem) FuncName() string {
	return this.Func().Name()
}

func (this *GoItem) Stop(wait ...bool) {
	this.StopTime = time.Now()
	this.Runner.Stop(wait...)
}

func (this *GoItem) run(ctx context.Context) error {
	this.StarTime = time.Now()
	this.StopTime = time.Time{}
	this.StoppedTime = time.Time{}
	this.f(ctx, this.Input...)
	this.StoppedTime = time.Now()
	return nil
}
