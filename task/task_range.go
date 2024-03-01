package task

import (
	"context"
	"errors"
	"github.com/injoyai/base/chans"
	"time"
)

type Handler func(ctx context.Context) (interface{}, error)

func NewRange() *Range {
	return &Range{
		limit: 1,
		retry: 3,
	}
}

type Range struct {
	queue    []Handler                                 //分片队列
	limit    uint                                      //协程数
	retry    uint                                      //重试次数
	doneItem func(ctx context.Context, resp *ItemResp) //子项执行完成
}

func (this *Range) Len() int {
	return len(this.queue)
}

func (this *Range) Append(v Handler) *Range {
	this.queue = append(this.queue, v)
	return this
}

func (this *Range) SetLimit(limit uint) *Range {
	this.limit = limit
	return this
}

func (this *Range) SetRetry(retry uint) *Range {
	this.retry = retry
	return this
}

func (this *Range) SetDoneItem(f func(ctx context.Context, resp *ItemResp)) *Range {
	this.doneItem = f
	return this
}

func (this *Range) Run(ctx context.Context) *Resp {
	start := time.Now()
	wg := chans.NewWaitLimit(this.limit)
	for i, f := range this.queue {
		select {
		case <-ctx.Done():
			return &Resp{Err: errors.New("上下文关闭")}
		default:
			wg.Add()
			go func(ctx context.Context, i int, f Handler) {
				defer wg.Done()
				resp := &ItemResp{Index: i}
				for x := uint(0); x < this.retry; x++ {
					resp.Data, resp.Err = f(ctx)
					if resp.Err == nil {
						break
					}
				}
				if this.doneItem != nil {
					this.doneItem(ctx, resp)
				}
			}(ctx, i, f)
		}
	}
	wg.Wait()
	return &Resp{
		Start: start,
	}
}

type ItemResp struct {
	Index int
	Data  interface{}
	Err   error
}

type Resp struct {
	Start time.Time     //任务开始时间
	Data  []interface{} //任务分片数据
	Err   error         //错误信息
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
