package task

import (
	"context"
	"errors"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/goutil/oss"
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
	queue     []GetBytes                                        //分片队列
	coroutine uint                                              //协程数
	retry     uint                                              //重试次数
	offset    int                                               //偏移量
	doneItem  func(ctx context.Context, resp *DownloadItemResp) //分片下载完成事件
}

func (this *Download) Len() int {
	return len(this.queue)
}

func (this *Download) Append(i int, v GetBytes) *Download {
	if v != nil {
		for len(this.queue) <= i {
			this.queue = append(this.queue, nil)
		}
		this.queue[i] = v
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

func (this *Download) SetDoneItem(doneItem func(ctx context.Context, resp *DownloadItemResp)) *Download {
	this.doneItem = doneItem
	return this
}

// Download 下载任务开始下载
func (this *Download) Download(ctx context.Context) *DownloadResp {
	if len(this.queue) == 1 {
		return this.downloadOne(ctx, this.queue[0])
	}
	start := time.Now()
	wg := chans.NewWaitLimit(this.coroutine)
	size := int64(0)
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
			go func(ctx context.Context, i int, t *Download, v GetBytes) {
				defer wg.Done()
				bytes, err := t.getBytes(ctx, v)
				if err == nil {
					size += int64(len(bytes))
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
		Size:  size,
	}
}

func (this *Download) downloadOne(ctx context.Context, i GetBytes) *DownloadResp {
	start := time.Now()
	bytes, err := i.GetBytes(ctx, func(p *http.Plan) {
		if this.doneItem != nil {
			this.doneItem(ctx, &DownloadItemResp{
				Index: p.Index,
				Bytes: p.Bytes,
			})
		}
	})
	return &DownloadResp{
		Start: start,
		Size:  int64(len(bytes)),
		Err:   err,
	}
}

func (this *Download) getBytes(ctx context.Context, v GetBytes) (bytes []byte, err error) {
	for i := uint(0); i < this.retry; i++ {
		bytes, err = v.GetBytes(ctx, func(p *http.Plan) {})
		if err == nil {
			return
		}
	}
	return
}

type GetBytes interface {
	GetBytes(ctx context.Context, f func(p *http.Plan)) ([]byte, error)
}

type GetBytesFunc func(ctx context.Context, f func(p *http.Plan)) ([]byte, error)

func (this GetBytesFunc) GetBytes(ctx context.Context, f func(p *http.Plan)) ([]byte, error) {
	return this(ctx, f)
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

func (this *DownloadItemResp) Save(filename string) error {
	return oss.New(filename, this.Bytes)
}

// DownloadResp 下载响应,去除字节,避免内存占用过大
type DownloadResp struct {
	Start time.Time      //任务开始时间
	Size  int64          //实际下载字节大小
	Err   error          //错误信息
	spend *time.Duration //下载耗时
	//Bytes [][]byte  //任务分片字节
}

//func (this *DownloadResp) WriteTo(w io.Writer) (int64, error) {
//	co := int64(0)
//	for _, bs := range this.Bytes {
//		if w != nil && bs != nil {
//			n, err := w.Write(bs)
//			if err != nil {
//				return co, err
//			}
//			co += int64(n)
//		}
//	}
//	if this.size == nil {
//		this.size = &co
//	}
//	return co, nil
//}
//
//func (this *DownloadResp) WriteToFile(filename string) (int64, error) {
//	f, err := os.Create(filename)
//	if err != nil {
//		return 0, err
//	}
//	defer f.Close()
//	return this.WriteTo(f)
//}

func (this *DownloadResp) GetSpend() time.Duration {
	if this.spend != nil {
		return *this.spend
	}
	spend := time.Since(this.Start)
	this.spend = &spend
	return spend
}

//func (this *DownloadResp) GetBytes() []byte {
//	bs := []byte(nil)
//	for _, v := range this.Bytes {
//		bs = append(bs, v...)
//	}
//	return bs
//}
//
//func (this *DownloadResp) GetSize() int64 {
//	if this.size != nil {
//		return *this.size
//	}
//	var size int64
//	for _, v := range this.Bytes {
//		size += int64(len(v))
//	}
//	this.size = &size
//	return size
//}

func (this *DownloadResp) Error() string {
	if this.Err != nil {
		return this.Err.Error()
	}
	return ""
}
