package http

import (
	"io"
	"time"
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
	return NewRequest(url)
}

// Get 使用默认客户端发起GET请求
func Get(url string, bind ...any) (*Response, error) {
	return DefaultClient.Get(url, bind...)
}

// GetBytes 使用GET请求获取响应字节
func GetBytes(url string) ([]byte, error) {
	return DefaultClient.GetBytes(url)
}

// GetReader 获取ReadCloser
func GetReader(url string) (io.ReadCloser, error) {
	return DefaultClient.GetReader(url)
}

// GetToWriter 使用GET请求获取响应字节写入writer,适用于下载请求
func GetToWriter(url string, writer io.Writer) (int64, error) {
	return DefaultClient.GetToWriter(url, writer)
}

// GetToFile 发起请求,并把body内容写入文件,适用于下载文件
func GetToFile(url string, filename string) (int64, error) {
	return DefaultClient.GetToFile(url, filename)
}

// Download 下载文件
func Download(url string, filename string) (int64, error) {
	return DefaultClient.GetToFile(url, filename)
}

func Post(url string, body any, bind ...any) (*Response, error) {
	return DefaultClient.Post(url, body, bind...)
}

func Put(url string, body any, bind ...any) (*Response, error) {
	return DefaultClient.Put(url, body, bind...)
}

func Delete(url string, body any, bind ...any) (*Response, error) {
	return DefaultClient.Delete(url, body, bind...)
}
