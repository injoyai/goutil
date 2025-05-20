package cache

import (
	"fmt"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
)

// Cycle 固定列表长度,循环使用
type Cycle struct {
	list       []any                 //列表数据
	offset     int                   //当前数据位置下标
	length     int                   //列表长度
	cycle      bool                  //循环使用,数据量已经超过列表的长度,覆盖了老数据
	subscribe  *chans.Subscribe[any] //数据订阅
	middleware []func(data any) bool //中间件
}

// Subscribe 开启一个订阅数据的通道
func (this *Cycle) Subscribe(cap ...uint) *chans.Safe[any] {
	return this.subscribe.Subscribe(cap...)
}

// Padding 填充数据,未测试
func (this *Cycle) Padding(data any) *Cycle {
	for i := range this.list {
		this.list[i] = data
	}
	this.cycle = true
	return this
}

// Save 数据持久化,保存至文件
func (this *Cycle) Save(name string) error {
	return newFile(name, "cycle").
		Set("data", this.list).
		Set("offset", this.offset).
		Set("length", this.length).
		Set("cycle", this.cycle).Save()
}

// List 获取列表数据(时间正序)
func (this *Cycle) List(limits ...int) []any {
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
func (this *Cycle) Use(f ...func(data any) bool) *Cycle {
	this.middleware = append(this.middleware, f...)
	return this
}

// Write 实现io.Writer接口
func (this *Cycle) Write(p []byte) (int, error) {
	this.Add(p)
	return len(p), nil
}

// Add 添加任意数据到缓存
func (this *Cycle) Add(data any) *Cycle {
	for _, f := range this.middleware {
		if !f(data) {
			return this
		}
	}
	this.subscribe.Publish(data)
	this.offset = conv.Select[int](this.offset >= len(this.list) || this.offset < 0, 0, this.offset)
	this.list[this.offset] = data
	this.offset++
	if this.offset >= this.length {
		this.offset = 0
		this.cycle = true
	}
	return this
}

func newCycle(length int) *Cycle {
	return &Cycle{
		list:      make([]any, length),
		offset:    0,
		length:    length,
		subscribe: chans.NewSubscribe[any](),
	}
}

// LoadingCycle 加载数据
func LoadingCycle(name string) (*Cycle, error) {
	f := newFile(name, "cycle")
	length := f.GetInt("length")
	if length <= 0 {
		return nil, fmt.Errorf("长度错误: %d", length)
	}
	c := newCycle(length)
	c.offset = f.GetInt("offset")
	c.cycle = f.GetBool("cycle")
	val, ok := f.MustGet("data").([]any)
	if ok {
		c.list = val
	}
	return c, nil
}
