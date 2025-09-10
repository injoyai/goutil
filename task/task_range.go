package task

import (
	"context"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"time"
)

type Handler[T any] func(ctx context.Context) (T, error)

func NewRange[T any]() *Range[T] {
	return &Range[T]{
		coroutine: 1,
		retry:     3,
	}
}

type Range[T any] struct {
	queue          []Handler[T]            //分片队列
	coroutine      int                     //协程数
	retry          int                     //重试次数
	onFinishedItem func(resp *ItemResp[T]) //子项执行完成
	onFinished     func(resp *Resp)        //全部完成
}

func (this *Range[T]) Len() int {
	num := 0
	for _, f := range this.queue {
		if f != nil {
			num++
		}
	}
	return num
}

func (this *Range[T]) Set(i int, v Handler[T]) *Range[T] {
	if v != nil {
		for len(this.queue) <= i {
			this.queue = append(this.queue, nil)
		}
		this.queue[i] = v
	}
	return this
}

func (this *Range[T]) Append(v ...Handler[T]) *Range[T] {
	this.queue = append(this.queue, v...)
	return this
}

func (this *Range[T]) SetCoroutine(limit int) *Range[T] {
	this.coroutine = conv.Select(limit == 0, 1, limit)
	return this
}

func (this *Range[T]) SetRetry(retry int) *Range[T] {
	this.retry = conv.Select(retry == 0, 1, retry)
	return this
}

func (this *Range[T]) OnFinishedItem(f func(resp *ItemResp[T])) *Range[T] {
	this.onFinishedItem = f
	return this
}

func (this *Range[T]) OnFinished(f func(resp *Resp)) *Range[T] {
	this.onFinished = f
	return this
}

func (this *Range[T]) Run(ctx context.Context) error {
	start := time.Now()
	wg := chans.NewWaitLimit(this.coroutine)
	for i, f := range this.queue {
		if f == nil {
			continue
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wg.Add()
			go func(ctx context.Context, t *Range[T], i int, f Handler[T]) {
				defer wg.Done()
				resp := &ItemResp[T]{Index: i}
				for x := 0; x <= t.retry; x++ {
					resp.Data, resp.Err = f(ctx)
					if resp.Err == nil {
						break
					}
				}
				if t.onFinishedItem != nil {
					t.onFinishedItem(resp)
				}
			}(ctx, this, i, f)
		}
	}
	wg.Wait()
	if this.onFinished != nil {
		this.onFinished(&Resp{Start: start})
	}
	return nil
}

type ItemResp[T any] struct {
	Index int
	Err   error
	Data  T
}

func (this *ItemResp[T]) Error() string {
	if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}

type Resp struct {
	Start time.Time      //任务开始时间
	spend *time.Duration //耗时
}

func (this *Resp) Spend() time.Duration {
	if this.spend != nil {
		return *this.spend
	}
	spend := time.Since(this.Start)
	this.spend = &spend
	return spend
}
