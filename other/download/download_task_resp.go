package download

import (
	"io"
	"os"
	"time"
)

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
