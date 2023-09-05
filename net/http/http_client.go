package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"time"
)

var DefaultClient = NewClient()

// NewClient
// 新建HTTP请求客户端
func NewClient() *Client {
	data := &Client{
		Client: &http.Client{
			Transport: &http.Transport{
				//连接结束后会直接关闭,
				//否则会加到连接池复用
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					//设置可以访问HTTPS
					InsecureSkipVerify: true,
				},
			},
			//设置连接超时时间,连接成功后无效
			//连接成功后数据读取时间可以超过这个时间
			//数据读取超时等可以nginx配置
			Timeout: time.Second * 10,
		},
	}
	return data
}

type Client struct {
	*http.Client
}

// SetProxy 设置代理
func (this *Client) SetProxy(u string) {
	if val, ok := this.Client.Transport.(*http.Transport); ok {
		val.Proxy = func(request *http.Request) (*url.URL, error) {
			if len(u) > 0 {
				return url.Parse(u)
			}
			return request.URL, nil
		}
	}
}

// SetTimeout 设置请求超时时间
func (this *Client) SetTimeout(t time.Duration) *Client {
	if this.Client != nil {
		this.Client.Timeout = t
	}
	return this
}

func (this *Client) Do(request *Request) (resp *Response) {
	start := time.Now()
	defer func() {
		if request.Done() {
			request.Reset()
			request.AddCookie(resp.Cookies()...)
		}
	}()
	request.AddTry()
	request.Request.Body = io.NopCloser(bytes.NewReader(request.body))
	r, err := this.Client.Do(request.Request)
	resp = newResponse(request, r, err).setStartTime(start)
	if resp.Err() != nil && !request.Done() {
		return this.Do(request)
	}
	return
}
