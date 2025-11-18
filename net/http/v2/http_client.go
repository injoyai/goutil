package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/proxy"
)

var DefaultClient = NewClient()

type Option func(c *Client)

func WithProxy(proxy string) Option {
	return func(c *Client) {
		c.SetProxy(proxy)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.SetTimeout(timeout)
	}
}

// NewClient
// 新建HTTP请求客户端
func NewClient(op ...Option) *Client {
	c := &Client{
		Client: &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: time.Second * 10,
		},
	}
	for _, v := range op {
		v(c)
	}
	return c
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

/*



 */

func (this *Client) Get(url string, bind ...any) (*Response, error) {
	resp, err := this.DoRequest(url, WithGet())
	if err != nil {
		return nil, err
	}
	if len(bind) > 0 {
		err = resp.Bind(bind[0])
	}
	return resp, err
}

func (this *Client) GetBytes(url string) ([]byte, error) {
	resp, err := this.DoRequest(url, WithGet())
	if err != nil {
		return nil, err
	}
	return resp.ReadBody()
}

func (this *Client) GetReader(url string) (io.ReadCloser, error) {
	resp, err := this.DoRequest(url, WithGet())
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

func (this *Client) GetToWriter(url string, w io.Writer) (int64, error) {
	resp, err := this.DoRequest(url, WithGet())
	if err != nil {
		return 0, err
	}
	defer resp.Response.Body.Close()
	return io.Copy(w, resp.Response.Body)
}

func (this *Client) GetToFile(url string, filename string) (int64, error) {
	resp, err := this.DoRequest(url, WithGet())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	w, err := os.Create(filename + ".downloading")
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(w, resp.Body)
	if err != nil {
		w.Close()
		return n, err
	}
	w.Close()
	<-time.After(time.Millisecond * 100)
	return n, os.Rename(filename+".downloading", filename)
}

func (this *Client) Post(url string, body any, bind ...any) (*Response, error) {
	resp, err := this.DoRequest(url, WithPost(), WithBody(body))
	if err != nil {
		return nil, err
	}
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp, nil
}

func (this *Client) Put(url string, body any, bind ...any) (*Response, error) {
	resp, err := this.DoRequest(url, WithPut(), WithBody(body))
	if err != nil {
		return nil, err
	}
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp, nil
}

func (this *Client) Delete(url string, body any, bind ...any) (*Response, error) {
	resp, err := this.DoRequest(url, WithDelete())
	if err != nil {
		return nil, err
	}
	if len(bind) > 0 {
		resp.Bind(bind[0])
	}
	return resp, nil
}

func (this *Client) Head(url string) (*Response, error) {
	return this.DoRequest(url, WithHead())
}

func (this *Client) Patch(url string) (*Response, error) {
	return this.DoRequest(url, WithPatch())
}

func (this *Client) Connect(url string) (*Response, error) {
	return this.DoRequest(url, WithConnect())
}

func (this *Client) Options(url string) (*Response, error) {
	return this.DoRequest(url, WithOptions())
}

func (this *Client) Trace(url string) (*Response, error) {
	return this.DoRequest(url, WithTrace())
}

func (this *Client) DoRequest(url string, op ...RequestOption) (*Response, error) {
	return this.Do(NewRequest(url, op...))
}

func (this *Client) Do(request *Request) (resp *Response, err error) {
	if request.Err() != nil {
		return nil, request.Err()
	}

	start := time.Now()
	defer func() {
		if resp != nil {
			request.AddCookie(resp.Cookies()...)
		}
	}()

	req, err := request.Request()
	if err != nil {
		return nil, err
	}
	var res *http.Response
	for i := 0; i == 0 || i < request.retry; i++ {
		res, err = this.Client.Do(req)
		if err == nil {
			resp = newResponse(res, start, i)
			if this.debug || request.debug {
				bs, _ := httputil.DumpRequest(req, true)
				fmt.Println("----------------------------------------------------------------")
				fmt.Println(string(bs))
				fmt.Println("----------------------------------------------------------------")
				fmt.Println(resp.String())
				fmt.Println("----------------------------------------------------------------")
			}
			return
		}
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
