package download

import (
	"context"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"io"
	"os"
	"time"
)

func NewTask() *Task {
	return &Task{
		limit: 1,
		retry: 1,
	}
}

type Task struct {
	queue    []GetBytes                                         //分片队列
	limit    uint                                               //协程数
	retry    uint                                               //重试次数
	offset   int                                                //偏移量
	doneItem func(ctx context.Context, resp *DoneItemResp)      //分片下载完成事件
	doneAll  func(ctx context.Context, resp *DoneAllResp) error //全部下载完成事件
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

func (this *Task) SetDoneAll(doneAll func(ctx context.Context, resp *DoneAllResp) error) *Task {
	this.doneAll = doneAll
	return this
}

func (this *Task) SetDoneAllWithWriter(w io.Writer) *Task {
	return this.SetDoneAll(func(ctx context.Context, resp *DoneAllResp) error {
		for _, bs := range resp.Bytes {
			if w != nil {
				if _, err := w.Write(bs); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (this *Task) SetDoneAllWithFile(filename string) *Task {
	return this.SetDoneAll(func(ctx context.Context, resp *DoneAllResp) error {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		for _, bs := range resp.Bytes {
			if _, err := f.Write(bs); err != nil {
				return err
			}
		}
		return nil
	})
}

// Download 下载任务开始下载
func (this *Task) Download(ctx context.Context) error {
	start := time.Now()
	wg := chans.NewWaitLimit(this.limit)
	cache := make([][]byte, this.Len())
	for i, v := range this.queue {
		select {
		case <-ctx.Done():
			return nil
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
	resp := &DoneAllResp{
		Start: start,
		Bytes: cache,
	}
	if this.doneAll != nil {
		return this.doneAll(ctx, resp)
	}
	return nil
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

type DoneItemResp struct {
	Index int
	Bytes []byte
	Err   error
}

func (this *DoneItemResp) GetSize() int64 {
	return int64(len(this.Bytes))
}

func (this *DoneItemResp) Error() string {
	if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}

type DoneAllResp struct {
	Start time.Time //任务开始时间
	Bytes [][]byte  //任务分片字节
	size  *int64
	spend *time.Duration
}

func (this *DoneAllResp) GetSpend() time.Duration {
	if this.spend != nil {
		return *this.spend
	}
	spend := time.Since(this.Start)
	this.spend = &spend
	return spend
}

func (this *DoneAllResp) GetBytes() []byte {
	bs := []byte(nil)
	for _, v := range this.Bytes {
		bs = append(bs, v...)
	}
	return bs
}

func (this *DoneAllResp) GetSize() int64 {
	if this.size != nil {
		return *this.size
	}
	var size int64
	for _, v := range this.Bytes {
		size += int64(len(v))
	}
	this.size = &size
	return size
}
