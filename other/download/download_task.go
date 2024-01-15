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

type TaskResp struct {
	Start time.Time //任务开始时间
	Bytes [][]byte  //任务分片字节
	size  *int64
	spend *time.Duration
}

func (this *TaskResp) WriteTo(w io.Writer) (int64, error) {
	co := int64(0)
	for _, bs := range this.Bytes {
		if w != nil && bs != nil {
			n, err := w.Write(bs)
			if err != nil {
				return co, err
			}
			co += int64(n)
		}
	}
	if this.size == nil {
		this.size = &co
	}
	return co, nil
}

func (this *TaskResp) WriteToFile(filename string) (int64, error) {
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return this.WriteTo(f)
}

func (this *TaskResp) GetSpend() time.Duration {
	if this.spend != nil {
		return *this.spend
	}
	spend := time.Since(this.Start)
	this.spend = &spend
	return spend
}

func (this *TaskResp) GetBytes() []byte {
	bs := []byte(nil)
	for _, v := range this.Bytes {
		bs = append(bs, v...)
	}
	return bs
}

func (this *TaskResp) GetSize() int64 {
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
