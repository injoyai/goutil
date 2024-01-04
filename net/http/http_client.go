package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"os"
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
func (this *Client) SetProxy(u string) *Client {
	if val, ok := this.Client.Transport.(*http.Transport); ok {
		uri, err := url.Parse(u)
		if err != nil || len(u) == 0 || len(uri.Host) == 0 {
			val.Proxy = nil
			return this
		}
		val.Proxy = func(request *http.Request) (*url.URL, error) {
			return uri, nil
		}
	}
	return this
}

// SetTimeout 设置请求超时时间
func (this *Client) SetTimeout(t time.Duration) *Client {
	if this.Client != nil {
		this.Client.Timeout = t
	}
	return this
}

func (this *Client) Get(url string, bind ...interface{}) *Response {
	resp := this.Do(NewRequest(http.MethodGet, url, nil))
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp
}

func (this *Client) GetBytes(url string) ([]byte, error) {
	resp := this.Do(NewRequest(http.MethodGet, url, nil))
	return resp.GetBodyBytes(), resp.Err()
}

func (this *Client) GetBytesWith(url string, f func([]byte)) ([]byte, error) {
	resp := this.Do(NewRequest(http.MethodGet, url, nil))
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	buf := bytes.NewBuffer(nil)
	if _, err := resp.CopyWith(buf, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (this *Client) GetBytesWithPlan(url string, f func(p *Plan)) ([]byte, error) {
	resp := this.Do(NewRequest(http.MethodGet, url, nil))
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	buf := bytes.NewBuffer(nil)
	if _, err := resp.CopyWithPlan(buf, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (this *Client) GetToWriter(url string, writer io.Writer) error {
	resp := this.Do(NewRequest(http.MethodGet, url, nil))
	if resp.Err() != nil {
		return resp.Err()
	}
	defer resp.Response.Body.Close()
	_, err := io.Copy(writer, resp.Response.Body)
	return err
}

func (this *Client) GetToFile(url string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return this.GetToWriter(url, f)
}

func (this *Client) Post(url string, body interface{}, bind ...interface{}) *Response {
	resp := this.Do(NewRequest(http.MethodPost, url, body))
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp
}

func (this *Client) Put(url string, body interface{}, bind ...interface{}) *Response {
	resp := this.Do(NewRequest(http.MethodPut, url, body))
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp
}

func (this *Client) Delete(url string, body interface{}, bind ...interface{}) *Response {
	resp := this.Do(NewRequest(http.MethodDelete, url, body))
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp
}

func (this *Client) Head(url string) *Response {
	return this.Do(NewRequest(http.MethodHead, url, nil))
}

func (this *Client) Patch(url string) *Response {
	return this.Do(NewRequest(http.MethodPatch, url, nil))
}

func (this *Client) Connect(url string) *Response {
	return this.Do(NewRequest(http.MethodConnect, url, nil))
}

func (this *Client) Options(url string) *Response {
	return this.Do(NewRequest(http.MethodOptions, url, nil))
}

func (this *Client) Trace(url string) *Response {
	return this.Do(NewRequest(http.MethodTrace, url, nil))
}

func (this *Client) Do(request *Request) (resp *Response) {
	start := time.Now()
	defer func() {
		if request.done() {
			request.reset()
			request.AddCookie(resp.Cookies()...)
		}
	}()
	request.addTry()
	request.Request.Body = io.NopCloser(bytes.NewReader(request.body))
	r, err := this.Client.Do(request.Request)
	resp = newResponse(request, r, err).setStartTime(start)
	if resp.Err() != nil && !request.done() {
		return this.Do(request)
	}
	return
}

func (this *Client) Forward(req *http.Request) *Response {
	return this.Do(&Request{
		Request: req,
		url:     req.URL.String(),
	})
}
