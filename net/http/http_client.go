package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/injoyai/io"
	"golang.org/x/net/proxy"
	"net"
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
	debug bool
}

func (this *Client) Debug(b ...bool) *Client {
	this.debug = len(b) == 0 || b[0]
	return this
}

// SetProxy 设置代理
func (this *Client) SetProxy(u string) error {
	if transport, ok := this.Client.Transport.(*http.Transport); ok {
		//为空表示取消代理
		if len(u) == 0 {
			transport.Proxy = nil
			transport.DialContext = nil
			return nil
		}
		proxyUrl, err := url.Parse(u)
		if err != nil {
			transport.Proxy = nil
			return err
		}
		switch proxyUrl.Scheme {
		case "socks5", "socks5h":
			dialer, err := proxy.FromURL(proxyUrl, this)
			if err != nil {
				return err
			}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		default: //"http", "https"
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
		return nil
	}
	return fmt.Errorf("http.Transport类型错误: 预期(*http.Transport),得到(%T)", this.Client.Transport)
}

// SetTimeout 设置请求超时时间
// 下载大文件的时候需要设置长的超时时间
func (this *Client) SetTimeout(t time.Duration) *Client {
	this.Client.Timeout = t
	return this
}

func (this *Client) Get(url string, bind ...interface{}) *Response {
	resp := this.DoRequest(http.MethodGet, url, nil)
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp
}

func (this *Client) GetBytes(url string) ([]byte, error) {
	resp := this.DoRequest(http.MethodGet, url, nil)
	return resp.GetBodyBytes(), resp.Err()
}

func (this *Client) GetReader(url string) (io.ReadCloser, error) {
	resp := this.DoRequest(http.MethodGet, url, nil)
	return resp.Body, resp.Err()
}

func (this *Client) GetToWriter(url string, w io.Writer) (int64, error) {
	resp := this.DoRequest(http.MethodGet, url, nil)
	if resp.Err() != nil {
		return 0, resp.Err()
	}
	defer resp.Response.Body.Close()
	return io.Copy(w, resp.Response.Body)
}

func (this *Client) GetToWriterWith(url string, w io.Writer, f func([]byte)) (int64, error) {
	resp := this.DoRequest(http.MethodGet, url, nil)
	if resp.Err() != nil {
		return 0, resp.Err()
	}
	return resp.CopyWith(w, f)
}

func (this *Client) GetToWriterWithPlan(url string, w io.Writer, f func(p *Plan)) (int64, error) {
	resp := this.DoRequest(http.MethodGet, url, nil)
	if resp.Err() != nil {
		return 0, resp.Err()
	}
	return resp.CopyWithPlan(w, func(p *io.Plan) {
		p.Total = resp.ContentLength
		f(p)
	})
}

func (this *Client) GetToFile(url string, filename string) (int64, error) {
	resp := this.DoRequest(http.MethodGet, url, nil)
	if resp.Err() != nil {
		return 0, resp.Err()
	}
	defer resp.Response.Body.Close()
	w, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer w.Close()
	return io.Copy(w, resp.Response.Body)
}

func (this *Client) GetToFileWithPlan(url string, filename string, f func(p *Plan)) (int64, error) {
	resp := this.DoRequest(http.MethodGet, url, nil)
	if resp.Err() != nil {
		return 0, resp.Err()
	}
	defer resp.Response.Body.Close()
	w, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer w.Close()
	return resp.CopyWithPlan(w, func(p *io.Plan) {
		p.Total = resp.ContentLength
		f(p)
	})
}

func (this *Client) Post(url string, body interface{}, bind ...interface{}) *Response {
	resp := this.DoRequest(http.MethodPost, url, body)
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp
}

func (this *Client) Put(url string, body interface{}, bind ...interface{}) *Response {
	resp := this.DoRequest(http.MethodPut, url, body)
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp
}

func (this *Client) Delete(url string, body interface{}, bind ...interface{}) *Response {
	resp := this.DoRequest(http.MethodDelete, url, body)
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp
}

func (this *Client) Head(url string) *Response {
	return this.DoRequest(http.MethodHead, url, nil)
}

func (this *Client) Patch(url string) *Response {
	return this.DoRequest(http.MethodPatch, url, nil)
}

func (this *Client) Connect(url string) *Response {
	return this.DoRequest(http.MethodConnect, url, nil)
}

func (this *Client) Options(url string) *Response {
	return this.DoRequest(http.MethodOptions, url, nil)
}

func (this *Client) Trace(url string) *Response {
	return this.DoRequest(http.MethodTrace, url, nil)
}

func (this *Client) DoRequest(method, url string, body interface{}) *Response {
	return this.Do(NewRequest(method, url, body))
}

func (this *Client) Do(request *Request) (resp *Response) {
	if request.Err() != nil {
		return newResponseErr(request.Err())
	}
	start := time.Now()
	defer func() {
		if request.done() {
			request.reset()
			request.AddCookie(resp.Cookies()...)
		}
	}()
	request.try++
	request.Request.Body = io.NopCloser(bytes.NewReader(request.body))
	r, err := this.Client.Do(request.Request)
	resp = newResponse(request, r, start, err)
	if this.debug || request.debug {
		fmt.Println(resp.String())
	}
	if resp.Err() != nil && !request.done() {
		return this.Do(request)
	}
	return
}

func (this *Client) Dial(network, addr string) (net.Conn, error) {
	d := &net.Dialer{
		Timeout:   this.Client.Timeout,
		KeepAlive: this.Client.Timeout,
	}
	return d.Dial(network, addr)
}

func (this *Client) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	d := &net.Dialer{
		Timeout:   this.Client.Timeout,
		KeepAlive: this.Client.Timeout,
	}
	return d.DialContext(ctx, network, addr)
}

/*



 */

func (this *Client) Request(url string, body ...interface{}) *Request {
	if len(body) > 0 {
		return NewRequest(http.MethodGet, url, body[0]).SetClient(this)
	}
	return NewRequest(http.MethodGet, url, nil).SetClient(this)
}

func (this *Client) Url(url string) *Request {
	return NewRequest(http.MethodGet, url, nil).SetClient(this)
}
