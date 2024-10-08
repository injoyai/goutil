package http

import (
	"bytes"
	"github.com/injoyai/conv"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	gourl "net/url"
	"sync"
)

type Request struct {
	*http.Request
	client *Client

	// query参数,?后面的参数
	// 可以可网站分开设置,所以另加这个字段用于保存参数
	// 例 xxx.SetQuery("a", 1).SetQuery("b", 2).SetUrl("http://www.baidu.com")
	// 结果为 http://www.baidu.com?a=1&b=2
	query   map[string]interface{}
	queryMu sync.RWMutex //锁

	// body参数,请求体,方便读取,提供GetBody等函数
	// 否则需要从流中读取,备份的body,当请求失败并需要重试时,可以从这里复制
	// todo 如果是文件的话,内存会占用比较大,如果优化
	body     []byte
	bodyBind interface{} //响应的body解析

	//debug模式,会打印请求响应的数据内容
	debug bool

	retry uint //重试次数,0不重试
	try   uint //已经执行次数,大于等于重试次数则结束

	muHeader sync.RWMutex //锁

	//执行中的错误信息,采用的链式操作,固先保存错误信息,统一处理
	err error
}

func (this *Request) Err() error {
	return this.err
}

// Reset 重置参数
// 采用指针类型,所以每次新请求参数都需要重新填写
// 加入预设值,重置之后运行预设方法
// 解决了客户端复用的问题(引用类型,复制客户端不一致问题)
func (this *Request) reset() {
	this.try = 0
}

// Done 判断是否完成,是否需要重试
func (this *Request) done() bool {
	return this.try > this.retry
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
	u, err := gourl.Parse(url)
	if err != nil {
		this.err = err
		return this
	}
	return this.SetURL(u)
}

func (this *Request) SetURL(u *gourl.URL) *Request {
	values := u.Query()
	this.queryMu.RLock()
	for i, v := range this.query {
		values.Add(i, conv.String(v))
	}
	this.queryMu.RUnlock()
	u.RawQuery = values.Encode()
	this.Request.URL = u
	this.Request.Host = u.Host
	return this
}

// GetUrl 获取请求地址
func (this *Request) GetUrl() string {
	if this.Request == nil || this.Request.URL == nil {
		return ""
	}
	return this.Request.URL.String()
}

// SetQuery 设置query参数,已存在则覆盖
func (this *Request) SetQuery(key string, val interface{}) *Request {
	this.queryMu.Lock()
	this.query[key] = val
	this.queryMu.Unlock()
	return this.SetURL(this.Request.URL)
}

// SetQuerys 批量设置query参数,已存在则覆盖
func (this *Request) SetQuerys(m map[string]interface{}) *Request {
	this.queryMu.Lock()
	for i, v := range m {
		this.query[i] = v
	}
	this.queryMu.Unlock()
	return this.SetURL(this.Request.URL)
}

// SetQueryMap 批量设置query参数,已存在则覆盖
func (this *Request) SetQueryMap(m map[string]interface{}) *Request {
	return this.SetQuerys(m)
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
func (this *Request) AddHeaders(m http.Header) *Request {
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
func (this *Request) SetHeaders(m http.Header) *Request {
	this.Request.Header = m
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
	this.SetContentType(w.FormDataContentType())
	return this.SetBody(body.Bytes())
}

// FormField form-data Field
func (this *Request) FormField(m map[string]interface{}) *Request {
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
func (this *Request) SetBody(i interface{}) *Request {
	this.body = conv.Bytes(i)
	this.Request.ContentLength = int64(len(this.body))
	this.Request.Body = io.NopCloser(bytes.NewReader(this.body))
	this.Request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(this.body)), nil
	}
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
func (this *Request) getClient() *Client {
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
	return this.getClient().Do(this)
}

// NewRequest 新建请求内容
func NewRequest(method, url string, body interface{}) *Request {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return &Request{
			// 消耗点内存,方便不用每个函数都进行不是nil的判断
			Request: &http.Request{Header: map[string][]string{}},
			err:     err,
		}
	}
	req := &Request{
		Request: request,
		query:   make(map[string]interface{}),
		err:     err,
	}
	req.SetBody(body)
	req.SetUserAgentDefault()
	return req
}
