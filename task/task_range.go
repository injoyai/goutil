package task

import (
	"context"
	"errors"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/goutil/g"
	"time"
)

func NewRange() *Range {
	return &Range{
		limit: 1,
		retry: 3,
	}
}

type Range struct {
	queue    []Handler              //分片队列
	limit    uint                   //协程数
	retry    uint                   //重试次数
	doneItem func(i int, err error) //子项执行完成
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

func (this *Range) SetDoneItem(f func(i int, err error)) *Range {
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
			go func(i int, f Handler) {
				defer wg.Done()
				err := g.Retry(f, int(this.retry))
				if this.doneItem != nil {
					this.doneItem(i, err)
				}
			}(i, f)
		}
	}
	wg.Wait()
	return &Resp{
		Start: start,
	}
}

type Resp struct {
	Start time.Time //任务开始时间
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
