package http

import (
	"bytes"
	"github.com/injoyai/conv"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	gourl "net/url"
	"strings"
	"sync"
)

type Request struct {
	*http.Request
	client   *Client
	url      string                 //网址
	query    map[string]interface{} //query参数,?后面的参数
	body     []byte                 //body参数
	bodyBind interface{}            //响应的body解析
	debug    bool                   //debug模式
	retry    uint                   //重试次数,0不重试
	try      uint                   //已经执行次数,大于等于重试次数则结束
	muQuery  sync.RWMutex           //锁
	muHeader sync.RWMutex           //锁
	err      error
}

// Reset 重置参数
// 采用指针类型,所以每次新请求参数都需要重新填写
// 加入预设值,重置之后运行预设方法
// 解决了客户端复用的问题(引用类型,复制客户端不一致问题)
func (this *Request) Reset() {
	this.try = 0
}

// Done 判断是否完成,是否需要重试
func (this *Request) Done() bool {
	return this.try > this.retry
}

// GetTry 获取重试次数
func (this *Request) GetTry() uint {
	return this.try
}

// AddTry 增加已重试次数
func (this *Request) AddTry() *Request {
	this.try++
	return this
}

// Retry 重试次数默认不重试
func (this *Request) Retry(num ...uint) *Request {
	if len(num) > 0 {
		this.retry = num[0]
	}
	return this
}

// Debug 打印请求响应参数
func (this *Request) Debug(debug ...bool) *Request {
	this.debug = !(len(debug) > 0 && !debug[0])
	return this
}

// String 输出字符串
func (this *Request) String() string {
	return string(this.Bytes())
}

// Bytes 转成字节,http协议,可直接通过tcp发送htp请求
func (this *Request) Bytes() []byte {
	bs, _ := httputil.DumpRequest(this.Request, true)
	bs = append(bs, this.body...)
	return bs
}

// SetMethod 设置请求方式
func (this *Request) SetMethod(method string) *Request {
	this.Request.Method = method
	return this
}

// GetMethod 获取请求方式
func (this *Request) GetMethod() string {
	return this.Request.Method
}

// SetUrl 设置地址
func (this *Request) SetUrl(url string) *Request {
	URL, err := gourl.Parse(this.dealQuery(url))
	if err == nil {
		this.url = url
		this.Request.URL = URL
	}
	return this
}

// GetUrl 获取请求地址
func (this *Request) GetUrl() string {
	return this.url
}

// SetQuery 设置query参数,已存在则覆盖
func (this *Request) SetQuery(key string, val interface{}) *Request {
	this.muQuery.Lock()
	this.query[key] = val
	this.muQuery.Unlock()
	return this.SetUrl(this.url)
}

// SetQueryMap 批量设置query参数,已存在则覆盖
func (this *Request) SetQueryMap(m map[string]interface{}) *Request {
	this.muQuery.Lock()
	for i, v := range m {
		this.query[i] = v
	}
	this.muQuery.Unlock()
	return this.SetUrl(this.url)
}

// copyHeader 复制请求头,map是引用类型
func (this *Request) copyHeader() map[string][]string {
	this.muHeader.RLock()
	defer this.muHeader.RUnlock()
	header := map[string][]string{}
	for i, v := range this.Request.Header {
		header[i] = v
	}
	return header
}

// AddHeader 添加请求头
func (this *Request) AddHeader(key string, val ...string) *Request {
	this.muHeader.Lock()
	defer this.muHeader.Unlock()
	this.Request.Header[key] = append(this.Request.Header[key], val...)
	return this
}

// AddHeaders 批量添加请求头header
func (this *Request) AddHeaders(m map[string][]string) *Request {
	for i, v := range m {
		this.AddHeader(i, v...)
	}
	return this
}

// SetHeader 设置请求头header,已存在则覆盖
func (this *Request) SetHeader(key string, val ...string) *Request {
	this.muHeader.Lock()
	defer this.muHeader.Unlock()
	this.Request.Header[key] = val
	return this
}

// SetHeaders 批量设置请求头header,,已存在则覆盖
func (this *Request) SetHeaders(m map[string][]string) *Request {
	header := http.Header{}
	for i, v := range m {
		header[i] = v
	}
	this.Request.Header = header
	return this
}

// AddCookie 添加请求头cookie
func (this *Request) AddCookie(cookies ...*http.Cookie) *Request {
	for _, cookie := range cookies {
		this.Request.AddCookie(cookie)
	}
	return this
}

// SetUserAgent 设置User-Agent
func (this *Request) SetUserAgent(s string) *Request {
	return this.SetHeader("User-Agent", s)
}

// SetUserAgentPostman 设置模拟成postman请求
func (this *Request) SetUserAgentPostman() *Request {
	return this.SetUserAgent("PostmanRuntime/7.35.0")
}

// SetUserAgentDefault 设置模拟成浏览器
func (this *Request) SetUserAgentDefault() *Request {
	return this.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
}

// SetReferer 设置Referer
func (this *Request) SetReferer(s string) *Request {
	return this.SetHeader("Referer", s)
}

// SetAuthorization 设置请求头Authorization
func (this *Request) SetAuthorization(s string) *Request {
	return this.SetHeader(HeaderAuthorization, s)
}

// SetToken 设置请求头Authorization,别名
func (this *Request) SetToken(s string) *Request {
	return this.SetAuthorization(s)
}

// SetContentType 设置请求头Content-Type
func (this *Request) SetContentType(s string) *Request {
	return this.SetHeader(HeaderContentType, s)
}

// FormFile form-data file
func (this *Request) FormFile(m map[string][]byte) *Request {
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
	this.SetContentType(w.FormDataContentType()).SetBody(body.Bytes())
	return this
}

// FormField form-data Field
func (this *Request) FormField(m map[string]interface{}) *Request {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for i, v := range m {
		_ = w.WriteField(i, conv.String(v))
	}
	_ = w.Close()
	this.SetContentType(w.FormDataContentType()).SetBody(body.Bytes())
	return this
}

// SetBody 设置请求body,默认json解析
func (this *Request) SetBody(i interface{}) *Request {
	this.body = conv.Bytes(i)
	this.Request.Body = io.NopCloser(bytes.NewReader(this.body))
	return this
}

// GetBody 获取body内容,返回字节
func (this *Request) GetBody() []byte {
	return this.body
}

// GetBodyString 获取body内容,返回字符串
func (this *Request) GetBodyString() string {
	return string(this.body)
}

// Bind 解析响应body,需要指针
func (this *Request) Bind(i interface{}) *Request {
	this.bodyBind = i
	return this
}

// SetClient 设置客户端
func (this *Request) SetClient(client *Client) *Request {
	if client != nil {
		this.client = client
	}
	return this
}

// GetClient 获取客户端
func (this *Request) GetClient() *Client {
	if this.client == nil {
		this.client = DefaultClient
	}
	if this.client == nil {
		this.client = NewClient()
	}
	return this.client
}

/*


 */

func (this *Request) GetBytes() ([]byte, error) {
	resp := this.Get()
	return resp.GetBody(), resp.Err()
}

func (this *Request) Get() *Response {
	return this.SetMethod(MethodGet).Do()
}

func (this *Request) Post() *Response {
	return this.SetMethod(MethodPost).Do()
}

func (this *Request) Put() *Response {
	return this.SetMethod(MethodPut).Do()
}

func (this *Request) Delete() *Response {
	return this.SetMethod(MethodDelete).Do()
}

func (this *Request) Do() *Response {
	if this.err != nil {
		return newResponse(nil, nil, this.err)
	}
	return this.GetClient().Do(this)
}

func (this *Request) dealQuery(url string) string {
	if len(this.query) > 0 {
		if !strings.Contains(url, "?") {
			url += "?"
		} else {
			url += "&"
		}
		u := gourl.Values{}
		this.muQuery.RLock()
		for i, v := range this.query {
			u.Add(i, conv.String(v))
		}
		this.muQuery.RUnlock()
		url += u.Encode()
	}
	return url
}

// NewRequest 新建请求内容
func NewRequest(method, url string, body interface{}) *Request {
	request, err := http.NewRequest(method, url, bytes.NewReader(conv.Bytes(body)))
	if request == nil {
		request = &http.Request{Header: map[string][]string{}}
	}
	data := &Request{
		Request: request,
		url:     url,
		query:   make(map[string]interface{}),
		body:    conv.Bytes(body),
		err:     err,
	}
	data.AddHeaders(map[string][]string{
		"User-Agent":   {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"},
		"Content-Type": {"application/json;charset=utf-8"}, //发送的数据格式
		"Accept":       {"application/json"},               //希望接收的数据格式
		"Connection":   {"close"},                          //短连接
	})
	return data
}
