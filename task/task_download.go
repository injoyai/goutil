package task

import (
	"context"
	"errors"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"io"
	"os"
	"time"
)

// NewDownload 新建下载任务
func NewDownload() *Download {
	return &Download{
		limit: 1,
		retry: 3,
	}
}

type Download struct {
	queue    []GetBytes                                        //分片队列
	limit    uint                                              //协程数
	retry    uint                                              //重试次数
	offset   int                                               //偏移量
	doneItem func(ctx context.Context, resp *DownloadItemResp) //分片下载完成事件
}

func (this *Download) Len() int {
	return len(this.queue)
}

func (this *Download) Append(v GetBytes) *Download {
	this.queue = append(this.queue, v)
	return this
}

func (this *Download) SetLimit(limit uint) *Download {
	this.limit = conv.SelectUint(limit == 0, 1, limit)
	return this
}

func (this *Download) SetRetry(retry uint) *Download {
	this.retry = conv.SelectUint(retry == 0, 1, retry)
	return this
}

func (this *Download) SetDoneItem(doneItem func(ctx context.Context, resp *DownloadItemResp)) *Download {
	this.doneItem = doneItem
	return this
}

// Download 下载任务开始下载
func (this *Download) Download(ctx context.Context) *DownloadResp {
	start := time.Now()
	wg := chans.NewWaitLimit(this.limit)
	cache := make([][]byte, this.Len())
	for i, v := range this.queue {
		select {
		case <-ctx.Done():
			return &DownloadResp{Err: errors.New("上下文关闭")}
		default:
			wg.Add()
			go func(ctx context.Context, i int, t *Download, v GetBytes) {
				defer wg.Done()
				bytes, err := t.getBytes(ctx, v)
				if err == nil {
					cache[i] = bytes
				}
				if this.doneItem != nil {
					this.doneItem(ctx, &DownloadItemResp{
						Index: i,
						Bytes: bytes,
						Err:   err,
					})
				}
			}(ctx, i, this, v)
		}
	}
	wg.Wait()
	return &DownloadResp{
		Start: start,
		Bytes: cache,
	}
}

func (this *Download) getBytes(ctx context.Context, v GetBytes) (bytes []byte, err error) {
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

/*



 */

type DownloadItemResp struct {
	Index int
	Bytes []byte
	Err   error
}

func (this *DownloadItemResp) GetSize() int64 {
	return int64(len(this.Bytes))
}

func (this *DownloadItemResp) Error() string {
	if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}

type DownloadResp struct {
	Start time.Time //任务开始时间
	Bytes [][]byte  //任务分片字节
	Err   error     //错误信息
	size  *int64
	spend *time.Duration
}

func (this *DownloadResp) WriteTo(w io.Writer) (int64, error) {
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

func (this *DownloadResp) WriteToFile(filename string) (int64, error) {
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return this.WriteTo(f)
}

func (this *DownloadResp) GetSpend() time.Duration {
	if this.spend != nil {
		return *this.spend
	}
	spend := time.Since(this.Start)
	this.spend = &spend
	return spend
}

func (this *DownloadResp) GetBytes() []byte {
	bs := []byte(nil)
	for _, v := range this.Bytes {
		bs = append(bs, v...)
	}
	return bs
}

func (this *DownloadResp) GetSize() int64 {
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

func (this *DownloadResp) Error() string {
	if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}
