package task

import (
	"context"
	"errors"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"io"
	"time"
)

// NewDownload 新建下载任务
func NewDownload() *Download {
	return &Download{
		coroutine: 1,
		retry:     3,
	}
}

type Download struct {
	queue     []GetReader                                                      //分片队列
	coroutine uint                                                             //协程数
	retry     uint                                                             //重试次数
	offset    int                                                              //偏移量
	doneItem  func(ctx context.Context, resp *DownloadItemResp) (int64, error) //分片下载完成事件
	doneAll   func(resp *DownloadResp)                                         //全部分片下载完成事件
}

func (this *Download) Len() int {
	return len(this.queue)
}

func (this *Download) Set(i int, v GetReader) *Download {
	if v != nil {
		for len(this.queue) <= i {
			this.queue = append(this.queue, nil)
		}
		this.queue[i] = v
	}
	return this
}

func (this *Download) Append(v GetReader) *Download {
	if v != nil {
		this.queue = append(this.queue, v)
	}
	return this
}

func (this *Download) SetCoroutine(limit uint) *Download {
	this.coroutine = conv.SelectUint(limit == 0, 1, limit)
	return this
}

func (this *Download) SetRetry(retry uint) *Download {
	this.retry = conv.SelectUint(retry == 0, 1, retry)
	return this
}

func (this *Download) SetDoneItem(doneItem func(ctx context.Context, resp *DownloadItemResp) (int64, error)) *Download {
	this.doneItem = doneItem
	return this
}

func (this *Download) SetDoneAll(doneAll func(resp *DownloadResp)) *Download {
	this.doneAll = doneAll
	return this
}

// Download 下载任务开始下载
func (this *Download) Download(ctx context.Context) *DownloadResp {
	if len(this.queue) == 1 {
		return this.downloadOne(ctx, this.queue[0])
	}
	start := time.Now()
	wg := chans.NewWaitLimit(this.coroutine)
	totalSize := int64(0)
	for i, v := range this.queue {
		if v == nil {
			continue
		}
		select {
		case <-ctx.Done():
			return &DownloadResp{
				Start: start,
				Size:  0,
				Err:   errors.New("上下文关闭"),
			}
		default:
			wg.Add()
			go func(ctx context.Context, i int, totalSize *int64, t *Download, v GetReader) {
				defer wg.Done()
				reader, err := t.getReader(ctx, v)
				if err == nil {
					defer reader.Close()
				}
				if t.doneItem == nil {
					t.doneItem = func(ctx context.Context, resp *DownloadItemResp) (int64, error) {
						return resp.GetSize()
					}
				}
				size, _ := t.doneItem(ctx, &DownloadItemResp{
					Index:  i,
					Err:    err,
					Reader: reader,
				})
				*totalSize += size
			}(ctx, i, &totalSize, this, v)
		}
	}
	wg.Wait()
	resp := &DownloadResp{
		Start: start,
		Size:  totalSize,
	}
	if this.doneAll != nil {
		this.doneAll(resp)
	}
	return resp
}

func (this *Download) downloadOne(ctx context.Context, i GetReader) *DownloadResp {
	start := time.Now()
	reader, err := this.getReader(ctx, i)
	size := int64(0)
	if this.doneItem != nil {
		size, _ = this.doneItem(ctx, &DownloadItemResp{
			Index:  0,
			Reader: reader,
		})
	}
	return &DownloadResp{
		Start: start,
		Size:  size,
		Err:   err,
	}
}

func (this *Download) getReader(ctx context.Context, v GetReader) (reader io.ReadCloser, err error) {
	for i := uint(0); i < this.retry; i++ {
		reader, err = v.GetReader(ctx)
		if err == nil {
			return
		}
	}
	return
}

type GetReader interface {
	GetReader(ctx context.Context) (io.ReadCloser, error)
}

type GetReaderFunc func(ctx context.Context) (io.ReadCloser, error)

func (this GetReaderFunc) GetReader(ctx context.Context) (io.ReadCloser, error) {
	return this(ctx)
}

/*



 */

type DownloadItemResp struct {
	Index  int
	Err    error
	Reader io.Reader
	bytes  *[]byte
}

func (this *DownloadItemResp) GetSize() (int64, error) {
	if this.bytes == nil && this.Reader != nil {
		bs, err := io.ReadAll(this.Reader)
		if err != nil {
			return 0, err
		}
		*this.bytes = bs
	}
	return int64(len(*this.bytes)), nil
}

func (this *DownloadItemResp) Error() string {
	if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}

func (this *DownloadItemResp) WriteTo(w io.Writer) (int64, error) {
	return io.Copy(w, this.Reader)
}

// Save 保存成文件
func (this *DownloadItemResp) Save(filename string) error {
	return oss.New(filename, this.Reader)
}

// DownloadResp 下载响应,去除字节,避免内存占用过大
type DownloadResp struct {
	Start time.Time      //任务开始时间
	Size  int64          //实际下载字节大小
	Err   error          //错误信息
	spend *time.Duration //下载耗时
}

func (this *DownloadResp) GetSpend() time.Duration {
	if this.spend != nil {
		return *this.spend
	}
	spend := time.Since(this.Start)
	this.spend = &spend
	return spend
}

func (this *DownloadResp) Error() string {
	if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}
