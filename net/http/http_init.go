package http

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/io"
	"net/http"
	"time"
)

const (
	MethodGet     = http.MethodGet
	MethodHead    = http.MethodHead
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodConnect = http.MethodConnect
	MethodOptions = http.MethodOptions
	MethodTrace   = http.MethodTrace
)

// SetProxy 设置默认客户端的代理地址
func SetProxy(proxy string) error {
	return DefaultClient.SetProxy(proxy)
}

// SetTimeout 设置超时时间
func SetTimeout(t time.Duration) *Client {
	return DefaultClient.SetTimeout(t)
}

func Url(url string) *Request {
	return NewRequest("", url, nil)
}

// Get 使用默认客户端发起GET请求
func Get(url string, bind ...interface{}) *Response {
	return DefaultClient.Get(url, bind...)
}

// GetReader 获取ReadCloser
func GetReader(url string) (io.ReadCloser, error) {
	return DefaultClient.GetReader(url)
}

// GetBytes 使用GET请求获取响应字节
func GetBytes(url string) ([]byte, error) {
	return DefaultClient.GetBytes(url)
}

// GetBodyDMap 获取body内容,解析成*conv.Map
func GetBodyDMap(url string) (*conv.Map, error) {
	resp := DefaultClient.Get(url)
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	return resp.GetBodyDMap(), nil
}

// GetToWriter 使用GET请求获取响应字节写入writer,适用于下载请求
func GetToWriter(url string, writer io.Writer) (int64, error) {
	return DefaultClient.GetToWriter(url, writer)
}

// GetToWriterWith 获取数据,分片复制数据
func GetToWriterWith(url string, w io.Writer, f func([]byte)) (int64, error) {
	return DefaultClient.GetToWriterWith(url, w, f)
}

// GetToWriterWithPlan 获取数据,分片复制数据
func GetToWriterWithPlan(url string, w io.Writer, fn func(p *Plan)) (int64, error) {
	return DefaultClient.GetToWriterWithPlan(url, w, fn)
}

// GetToFile 发起请求,并把body内容写入文件,适用于下载文件
func GetToFile(url string, filename string) (int64, error) {
	return DefaultClient.GetToFile(url, filename)
}

// Download 下载文件
func Download(url string, filename string, f ...func(p *Plan)) (int64, error) {
	if len(f) > 0 && f[0] != nil {
		return DefaultClient.GetToFileWithPlan(url, filename, f[0])
	}
	return DefaultClient.GetToFile(url, filename)
}

func Post(url string, body interface{}, bind ...interface{}) *Response {
	return DefaultClient.Post(url, body, bind...)
}

func Put(url string, body interface{}, bind ...interface{}) *Response {
	return DefaultClient.Put(url, body, bind...)
}

func Delete(url string, body interface{}, bind ...interface{}) *Response {
	return DefaultClient.Delete(url, body, bind...)
}
