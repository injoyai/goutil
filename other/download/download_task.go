package download

import (
	"context"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"time"
)

func NewTask() *Task {
	return &Task{
		limit: 1,
		retry: 1,
	}
}

type Task struct {
	queue    []GetBytes                                    //分片队列
	limit    uint                                          //协程数
	retry    uint                                          //重试次数
	offset   int                                           //偏移量
	doneItem func(ctx context.Context, resp *DoneItemResp) //分片下载完成事件
}

func (this *Task) Len() int {
	return len(this.queue)
}

func (this *Task) Append(v GetBytes) *Task {
	this.queue = append(this.queue, v)
	return this
}

func (this *Task) SetLimit(limit uint) *Task {
	this.limit = conv.SelectUint(limit == 0, 1, limit)
	return this
}

func (this *Task) SetRetry(retry uint) *Task {
	this.retry = conv.SelectUint(retry == 0, 1, retry)
	return this
}

func (this *Task) SetDoneItem(doneItem func(ctx context.Context, resp *DoneItemResp)) *Task {
	this.doneItem = doneItem
	return this
}

// Download 下载任务开始下载
func (this *Task) Download(ctx context.Context) *TaskResp {
	start := time.Now()
	wg := chans.NewWaitLimit(this.limit)
	cache := make([][]byte, this.Len())
	for i, v := range this.queue {
		select {
		case <-ctx.Done():
			return &TaskResp{}
		default:
			wg.Add()
			go func(ctx context.Context, i int, t *Task, v GetBytes) {
				defer wg.Done()
				bytes, err := t.getBytes(ctx, v)
				if err == nil {
					cache[i] = bytes
				}
				if this.doneItem != nil {
					this.doneItem(ctx, &DoneItemResp{
						Index: i,
						Bytes: bytes,
						Err:   err,
					})
				}
			}(ctx, i, this, v)
		}
	}
	wg.Wait()
	return &TaskResp{
		Start: start,
		Bytes: cache,
	}
}

func (this *Task) getBytes(ctx context.Context, v GetBytes) (bytes []byte, err error) {
	for i := uint(0); i < this.retry; i++ {
		bytes, err = v.GetBytes(ctx)
		if err == nil {
			return
		}
	}
	return
}

type GetBytes interface {
	GetBytes(ctx context.Context) ([]byte, error)
}

type GetBytesFunc func(ctx context.Context) ([]byte, error)

func (this GetBytesFunc) GetBytes(ctx context.Context) ([]byte, error) {
	return this(ctx)
}
