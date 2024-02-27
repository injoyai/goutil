package http

import (
	"github.com/injoyai/io"
	"net/http"
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

func Url(url string) *Request {
	return NewRequest("", url, nil)
}

// Get 使用默认客户端发起GET请求
func Get(url string, bind ...interface{}) *Response {
	return DefaultClient.Get(url, bind...)
}

// GetBytes 使用GET请求获取响应字节
func GetBytes(url string) ([]byte, error) {
	return DefaultClient.GetBytes(url)
}

// GetWith 获取数据并监听
func GetWith(url string, f func([]byte)) (int64, error) {
	return DefaultClient.GetWith(url, f)
}

// GetWithPlan 获取数据并监听
func GetWithPlan(url string, fn func(p *io.Plan)) (int64, error) {
	return DefaultClient.GetWithPlan(url, fn)
}

// GetToWriter 使用GET请求获取响应字节写入writer,适用于下载请求
func GetToWriter(url string, writer io.Writer) error {
	return DefaultClient.GetToWriter(url, writer)
}

// GetToFile 发起请求,并把body内容写入文件,适用于下载文件
func GetToFile(url string, filename string) error {
	return DefaultClient.GetToFile(url, filename)
}

// Download 下载文件
func Download(url string, filename string) error {
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
