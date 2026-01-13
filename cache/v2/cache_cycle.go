package cache

import (
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/types"
	"github.com/injoyai/conv"
)

func NewCycle[T any](length int) *Cycle[T] {
	return &Cycle[T]{
		list:      make([]T, length),
		offset:    0,
		length:    length,
		subscribe: chans.NewSubscribe[string, T](),
	}
}

// Cycle 固定列表长度,循环使用
type Cycle[T any] struct {
	list       types.List[T]               //列表数据
	offset     int                         //当前数据位置下标
	length     int                         //列表长度
	cycle      bool                        //循环使用,数据量已经超过列表的长度,覆盖了老数据
	subscribe  *chans.Subscribe[string, T] //数据订阅
	middleware []func(data T) bool         //中间件
}

// Subscribe 开启一个订阅数据的通道
func (this *Cycle[T]) Subscribe(cap ...int) *chans.Safe[T] {
	return this.subscribe.Subscribe("", cap...)
}

// Padding 填充数据
func (this *Cycle[T]) Padding(data T) *Cycle[T] {
	for i := range this.list {
		this.list[i] = data
	}
	this.cycle = true
	return this
}

// Save 数据持久化,保存至文件
func (this *Cycle[T]) Save(name string) error {
	return NewFile(name, "cycle").
		Set("data", this.list).
		Set("offset", this.offset).
		Set("length", this.length).
		Set("cycle", this.cycle).
		Save()
}

// List 获取列表数据(时间正序)
func (this *Cycle[T]) List(limits ...int) []T {
	this.offset = conv.Select[int](this.offset >= len(this.list) || this.offset < 0, 0, this.offset)
	list := this.list[:this.offset]
	if this.cycle {
		list = append(this.list[this.offset:], list...)
	}
	if len(limits) > 0 && len(list) > limits[0] {
		return list[len(list)-limits[0]:]
	}
	return list
}

// Use 中间件
func (this *Cycle[T]) Use(f ...func(data T) bool) *Cycle[T] {
	this.middleware = append(this.middleware, f...)
	return this
}

// Append 添加任意数据到缓存
func (this *Cycle[T]) Append(data T) *Cycle[T] {
	for _, f := range this.middleware {
		if !f(data) {
			return this
		}
	}
	this.subscribe.Publish("", data)
	this.offset = conv.Select[int](this.offset >= len(this.list) || this.offset < 0, 0, this.offset)
	this.list[this.offset] = data
	this.offset++
	if this.offset >= this.length {
		this.offset = 0
		this.cycle = true
	}
	return this
}
