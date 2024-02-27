package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/injoyai/conv"
	"github.com/injoyai/io"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type Response struct {
	*http.Response
	Request *Request
	body    []byte            //body数据
	spend   time.Duration     //花费时间
	tryNum  uint              //尝试次数
	Error   error             //错误信息
	doc     *goquery.Document //解析内容
}

func (this *Response) setErr(err error) {
	if err != nil && this.Error == nil {
		this.Error = err
	}
}

// String 实现系统接口,默认输出
func (this *Response) String() string {
	if this == nil || this.Response == nil {
		return ""
	}
	if this.Error != nil {
		return this.Error.Error()
	}
	if this.Request != nil && this.Response != nil {
		this.Response.Header.Add(HeaderKeySpend, this.spend.String())
		this.Response.Header.Add(HeaderKeyTry, conv.String(this.tryNum))
		respBs, err := httputil.DumpResponse(this.Response, true)
		if err != nil {
			respBs, _ = httputil.DumpResponse(this.Response, false)
		}
		/*
			http: ContentLength=7356416 with Body length 4701759
		*/
		return fmt.Sprintf(`----------------------------------------
%s 
----------------------------------------
%s
----------------------------------------
`, this.Request, string(respBs))
	}
	return ""
}

// TryNum 执行的次数,正常请求是1次
func (this *Response) TryNum() uint {
	return this.tryNum
}

// Spend 获取花费时间
func (this *Response) Spend() time.Duration {
	return this.spend
}

// Code 状态码
func (this *Response) Code() int {
	if this.Response == nil {
		return 500
	}
	return this.StatusCode
}

// Header 请求头
func (this *Response) Header() http.Header {
	if this.Response == nil {
		return http.Header{}
	}
	return this.Response.Header
}

// GetHeader 获取请求头
func (this *Response) GetHeader(key string) string {
	return this.Header().Get(key)
}

// GetContentLength 获取Content-Length
func (this *Response) GetContentLength() int64 {
	return this.ContentLength
}

// Cookies 获取cookie信息
func (this *Response) Cookies() (cookie []*http.Cookie) {
	if this.Response != nil {
		return this.Response.Cookies()
	}
	return
}

// CopyWith 复制数据,并监听
func (this *Response) CopyWith(w io.Writer, fn func(bs []byte)) (int64, error) {
	return io.CopyWith(w, this.Body, fn)
}

// CopyWithPlan 复制数据并监听进度
func (this *Response) CopyWithPlan(w io.Writer, fn func(p *io.Plan)) (int64, error) {
	return io.CopyWithPlan(w, this.Body, func(p *io.Plan) {
		p.Total = this.GetContentLength()
		fn(p)
	})
}

// WriteToFile 写入新文件,会覆盖原文件(如果存在)
func (this *Response) WriteToFile(filename string) (int64, error) {
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return this.WriteTo(f)
}

// WriteTo 写入到writer,例如文件下载,写入到文件
func (this *Response) WriteTo(writer io.Writer) (int64, error) {
	return io.Copy(writer, this.Response.Body)
}

// GetReadCloser 获取结果的body reader
func (this *Response) GetReadCloser() io.ReadCloser {
	return this.Response.Body
}

// GetBody 一次性读取全部字节(适用于小数据)
func (this *Response) GetBody() []byte {
	if this.body == nil && this.Response != nil {
		this.body, _ = io.ReadAll(this.Response.Body)
		this.Response.Body.Close()
		//保持body可读
		this.Body = io.NopCloser(bytes.NewBuffer(this.body))
	}
	return this.body
}

// GetBodyBytes 获取body内容,,返回字节
func (this *Response) GetBodyBytes() []byte {
	return this.GetBody()
}

// GetBodyString 获取body内容,返回字符串
func (this *Response) GetBodyString() string {
	return string(this.GetBody())
}

// GetBodyDMap 获取body内容,解析成*conv.Map
func (this *Response) GetBodyDMap() *conv.Map {
	return conv.NewMap(this.GetBodyBytes())
}

// GetBodyMap 获取body内容,解析成map[string]interface{}返回
func (this *Response) GetBodyMap() (m map[string]interface{}) {
	_ = json.Unmarshal(this.GetBodyBytes(), &m)
	return
}

// GetBodyMaps 获取body内容,解析成[]map[string]interface{}返回
func (this *Response) GetBodyMaps() (m []map[string]interface{}) {
	_ = json.Unmarshal(this.GetBodyBytes(), &m)
	return
}

// Execute 执行,方便操作
func (this *Response) Execute(f ...func(r *Response) error) (err error) {
	if this.Error != nil {
		return this.Error
	}
	for _, v := range f {
		if err := v(this); err != nil {
			return err
		}
	}
	return nil
}

// Bind 绑定body数据,目前支持字符串,字节和json,需要指针
func (this *Response) Bind(ptr interface{}) *Response {
	if ptr != nil {
		switch val := ptr.(type) {
		case *string:
			*val = this.GetBodyString()
		case *[]byte:
			//val不为nil,this.body为nil可以赋值成功
			*val = this.GetBodyBytes()
		case io.Writer:
			_, err := io.Copy(val, this.Body)
			this.setErr(err)
		default:
			this.setErr(json.Unmarshal(this.GetBodyBytes(), ptr))
		}
	}
	return this
}

// Err 错误信息,如果错误则Response为nil
func (this *Response) Err() error {
	if this.Error != nil {
		return this.Error
	}
	if this.Response != nil && this.Response.StatusCode/100 != 2 {
		return fmt.Errorf("状态码错误:%d\n%s", this.Response.StatusCode, this.GetBodyString())
	}
	return nil
}

func newResponseErr(err error) *Response {
	return &Response{Error: err}
}

func newResponse(req *Request, resp *http.Response, start time.Time, err error) *Response {
	r := &Response{
		Request:  req,
		Response: resp,
		spend:    time.Now().Sub(start),
		Error:    err,
	}
	if req != nil {
		r.tryNum = req.try
		r.Bind(req.bodyBind)
	}
	return r
}
