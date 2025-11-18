package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/injoyai/conv"
)

type RequestOption func(r *Request)

func WithMethod(method string) RequestOption {
	return func(r *Request) {
		r.SetMethod(method)
	}
}

func WithGet() RequestOption {
	return func(r *Request) {
		r.SetMethod(http.MethodGet)
	}
}

func WithPost() RequestOption {
	return func(r *Request) {
		r.SetMethod(http.MethodPost)
	}
}

func WithPut() RequestOption {
	return func(r *Request) {
		r.SetMethod(http.MethodPut)
	}
}

func WithDelete() RequestOption {
	return func(r *Request) {
		r.SetMethod(http.MethodDelete)
	}
}

func WithHead() RequestOption {
	return func(r *Request) {
		r.SetMethod(http.MethodHead)
	}
}

func WithOptions() RequestOption {
	return func(r *Request) {
		r.SetMethod(http.MethodOptions)
	}
}

func WithPatch() RequestOption {
	return func(r *Request) {
		r.SetMethod(http.MethodPatch)
	}
}

func WithTrace() RequestOption {
	return func(r *Request) {
		r.SetMethod(http.MethodTrace)
	}
}

func WithConnect() RequestOption {
	return func(r *Request) {
		r.SetMethod(http.MethodConnect)
	}
}

func WithUrl(url string) RequestOption {
	return func(r *Request) {
		r.SetUrl(url)
	}
}

func WithBody(body any) RequestOption {
	return func(r *Request) {
		r.SetBody(body)
	}
}

func WithHeader(k string, v string) RequestOption {
	return func(r *Request) {
		r.SetHeader(k, v)
	}
}

func WithHeaders(headers http.Header) RequestOption {
	return func(r *Request) {
		r.SetHeaders(headers)
	}
}

func WithQuery(k string, v any) RequestOption {
	return func(r *Request) {
		r.SetQuery(k, v)
	}
}

func WithQueries(query map[string]any) RequestOption {
	return func(r *Request) {
		r.SetQueries(query)
	}
}

func WithCookies(cookie ...*http.Cookie) RequestOption {
	return func(r *Request) {
		r.AddCookie(cookie...)
	}
}

// NewRequest 新建请求内容
func NewRequest(url string, op ...RequestOption) *Request {
	req := &Request{
		method: http.MethodGet,
		url:    url,
		query:  make(map[string]any),
		header: make(http.Header),
	}
	req.SetUserAgentDefault()
	for _, o := range op {
		o(req)
	}
	return req
}

func Url(url string) *Request {
	return NewRequest(url)
}

type Request struct {

	// 请求方式
	method string

	// 请求地址
	url string

	// query参数,?后面的参数
	// 可以可网站分开设置,所以另加这个字段用于保存参数
	// 例 xxx.SetQuery("a", 1).SetQuery("b", 2).SetUrl("http://www.baidu.com")
	// 结果为 http://www.baidu.com?a=1&b=2
	query map[string]any

	header http.Header

	cookies []*http.Cookie

	mu sync.RWMutex

	body any

	bind any //响应的body解析

	//debug模式,会打印请求响应的数据内容
	debug bool

	retry int //重试次数,0不重试

	//执行中的错误信息,采用的链式操作,固先保存错误信息,统一处理
	err error
}

func (this *Request) Request() (*http.Request, error) {
	return http.NewRequest(this.method, this.url, this.Body())
}

func (this *Request) Err() error {
	return this.err
}

// Retry 重试次数默认不重试
func (this *Request) Retry(num ...int) *Request {
	this.retry = conv.Default(0, num...)
	return this
}

// Debug 打印请求响应参数
func (this *Request) Debug(debug ...bool) *Request {
	this.debug = len(debug) == 0 || debug[0]
	return this
}

// String 输出字符串
func (this *Request) String() string {
	return string(this.Bytes())
}

// Bytes 转成字节,http协议,可直接通过tcp发送htp请求
func (this *Request) Bytes() []byte {
	req, err := http.NewRequest(this.method, this.url, this.Body())
	if err != nil {
		return nil
	}
	bs, _ := httputil.DumpRequest(req, true)
	return bs
}

// SetMethod 设置请求方式
func (this *Request) SetMethod(method string) *Request {
	this.method = method
	return this
}

// Method 获取请求方式
func (this *Request) Method() string {
	return this.method
}

// SetUrl 设置地址
func (this *Request) SetUrl(url string) *Request {
	this.url = url
	return this
}

// Url 获取请求地址
func (this *Request) Url() string {
	return this.url
}

// SetQuery 设置query参数,已存在则覆盖
func (this *Request) SetQuery(key string, val any) *Request {
	this.mu.Lock()
	this.query[key] = val
	this.mu.Unlock()
	return this
}

// SetQueries 批量设置query参数,已存在则覆盖
func (this *Request) SetQueries(m map[string]any) *Request {
	this.mu.Lock()
	this.query = m
	this.mu.Unlock()
	return this
}

// AddHeader 添加请求头
func (this *Request) AddHeader(key string, val ...string) *Request {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.header[key] = append(this.header[key], val...)
	return this
}

// AddHeaders 批量添加请求头header
func (this *Request) AddHeaders(m http.Header) *Request {
	for i, v := range m {
		this.AddHeader(i, v...)
	}
	return this
}

// SetHeader 设置请求头header,已存在则覆盖
func (this *Request) SetHeader(key string, val ...string) *Request {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.header[key] = val
	return this
}

// SetHeaders 批量设置请求头header,,已存在则覆盖
func (this *Request) SetHeaders(m http.Header) *Request {
	this.header = m
	return this
}

// AddCookie 添加请求头cookie
func (this *Request) AddCookie(cookies ...*http.Cookie) *Request {
	this.cookies = append(this.cookies, cookies...)
	return this
}

func (this *Request) SetCookie(cookies ...*http.Cookie) *Request {
	this.cookies = cookies
	return this
}

// SetUserAgent 设置User-Agent
func (this *Request) SetUserAgent(s string) *Request {
	return this.SetHeader(HeaderKeyUserAgent, s)
}

// SetUserAgentPostman 设置模拟成postman请求
func (this *Request) SetUserAgentPostman() *Request {
	return this.SetUserAgent(UserAgentPostman)
}

// SetUserAgentDefault 设置模拟成浏览器
func (this *Request) SetUserAgentDefault() *Request {
	return this.SetUserAgent(UserAgentDefault)
}

// SetReferer 设置Referer
func (this *Request) SetReferer(s string) *Request {
	return this.SetHeader(HeaderKeyReferer, s)
}

// SetAuthorization 设置请求头Authorization
func (this *Request) SetAuthorization(s string) *Request {
	return this.SetHeader(HeaderKeyAuthorization, s)
}

// SetToken 设置请求头Authorization,别名
func (this *Request) SetToken(s string) *Request {
	return this.SetAuthorization(s)
}

// SetContentType 设置请求头Content-Type
func (this *Request) SetContentType(s string) *Request {
	return this.SetHeader(HeaderKeyContentType, s)
}

// SetFormFile form-data file
func (this *Request) SetFormFile(m map[string][]byte) *Request {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for i, v := range m {
		fileWriter, err := w.CreateFormFile(i, i)
		if err != nil {
			continue
		}
		fileWriter.Write(v)
	}
	_ = w.Close()
	this.SetContentType(w.FormDataContentType())
	return this.SetBody(body.Bytes())
}

// SetFormField form-data Field
func (this *Request) SetFormField(m map[string]any) *Request {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for i, v := range m {
		_ = w.WriteField(i, conv.String(v))
	}
	_ = w.Close()
	this.SetContentType(w.FormDataContentType())
	return this.SetBody(body.Bytes())
}

// SetBody 设置请求body,默认json解析
func (this *Request) SetBody(i any) *Request {
	this.body = i
	return this
}

// Body 获取body内容,返回字节
func (this *Request) Body() io.Reader {
	switch r := this.body.(type) {
	case nil:
		return nil
	case []byte:
		return bytes.NewReader(r)
	case string:
		return strings.NewReader(r)
	case io.Reader:
		return r
	default:
		return strings.NewReader(conv.String(this.body))
	}
}

// Bind 解析响应body,需要指针
func (this *Request) Bind(i any) *Request {
	this.bind = i
	return this
}

/*



 */

func (this *Request) Get(c ...*Client) (*Response, error) {
	return this.SetMethod(http.MethodGet).Do(c...)
}

func (this *Request) Post(c ...*Client) (*Response, error) {
	return this.SetMethod(http.MethodPost).Do(c...)
}

func (this *Request) Put(c ...*Client) (*Response, error) {
	return this.SetMethod(http.MethodPut).Do(c...)
}

func (this *Request) Delete(c ...*Client) (*Response, error) {
	return this.SetMethod(http.MethodDelete).Do(c...)
}

func (this *Request) Do(c ...*Client) (*Response, error) {
	cli := conv.Default(DefaultClient, c...)
	resp, err := cli.Do(this)
	if err != nil {
		return nil, err
	}
	err = resp.Bind(this.bind)
	return resp, err
}
