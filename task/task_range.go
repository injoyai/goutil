package task

import (
	"context"
	"errors"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"time"
)

type Handler func(ctx context.Context) (any, error)

func NewRange() *Range {
	return &Range{
		coroutine: 1,
		retry:     3,
	}
}

type Range struct {
	queue     []Handler                                 //分片队列
	coroutine uint                                      //协程数
	retry     uint                                      //重试次数
	doneItem  func(ctx context.Context, resp *ItemResp) //子项执行完成
	doneAll   func(resp *Resp)                          //全部完成
}

func (this *Range) Len() int {
	return len(this.queue)
}

func (this *Range) Set(i int, v Handler) *Range {
	if v != nil {
		for len(this.queue) <= i {
			this.queue = append(this.queue, nil)
		}
		this.queue[i] = v
	}
	return this
}

func (this *Range) Append(v Handler) *Range {
	this.queue = append(this.queue, v)
	return this
}

func (this *Range) SetCoroutine(limit uint) *Range {
	this.coroutine = conv.Select[uint](limit == 0, 1, limit)
	return this
}

func (this *Range) SetRetry(retry uint) *Range {
	this.retry = conv.Select[uint](retry == 0, 1, retry)
	return this
}

func (this *Range) SetDoneItem(f func(ctx context.Context, resp *ItemResp)) *Range {
	this.doneItem = f
	return this
}

func (this *Range) SetDoneAll(doneAll func(resp *Resp)) *Range {
	this.doneAll = doneAll
	return this
}

func (this *Range) Run(ctx context.Context) *Resp {
	start := time.Now()
	wg := chans.NewWaitLimit(this.coroutine)
	for i, f := range this.queue {
		if f == nil {
			continue
		}
		select {
		case <-ctx.Done():
			return &Resp{Err: errors.New("上下文关闭")}
		default:
			wg.Add()
			go func(ctx context.Context, t *Range, i int, f Handler) {
				defer wg.Done()
				resp := &ItemResp{Index: i}
				for x := uint(0); x <= t.retry; x++ {
					resp.Data, resp.Err = f(ctx)
					if resp.Err == nil {
						break
					}
				}
				if t.doneItem != nil {
					t.doneItem(ctx, resp)
				}
			}(ctx, this, i, f)
		}
	}
	wg.Wait()
	resp := &Resp{
		Start: start,
	}
	if this.doneAll != nil {
		this.doneAll(resp)
	}
	return resp
}

type ItemResp struct {
	Index int
	Err   error
	Data  any
}

func (this *ItemResp) Error() string {
	if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}

type Resp struct {
	Start time.Time //任务开始时间
	Data  []any     //任务分片数据
	Err   error     //错误信息
	spend *time.Duration
}

func (this *Resp) GetSpend() time.Duration {
	if this.spend != nil {
		return *this.spend
	}
	spend := time.Since(this.Start)
	this.spend = &spend
	return spend
}

func (this *Resp) Error() string {
	if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}
