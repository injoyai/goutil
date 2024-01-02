package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type Response struct {
	*http.Response
	Request *Request
	body    []byte        //body数据
	spend   time.Duration //花费时间
	tryNum  uint          //尝试次数
	err     error         //错误信息
}

// print 打印输出信息
func (this *Response) print() *Response {
	if this.Request != nil && this.Request.debug && this.Response != nil {
		fmt.Print(this.String())
	}
	return this
}

// setTryNum 设置已重试的次数
func (this *Response) setTryNum(num uint) *Response {
	this.tryNum = num
	return this
}

// setStartTime 设置已花费的时间
func (this *Response) setStartTime(start time.Time) *Response {
	this.spend = time.Now().Sub(start)
	return this
}

// setErr 设置错误
func (this *Response) setErr(err error) {
	if err != nil && this.err == nil {
		this.err = err
	}
}

/*


 */

// String 实现系统接口,默认输出
func (this *Response) String() string {
	if this.err != nil {
		return this.err.Error()
	}
	if this.Request != nil && this.Response != nil {
		this.Response.Header.Add(HeaderKeySpend, this.spend.String())
		this.Response.Header.Add(HeaderKeyTry, conv.String(this.tryNum))
		respBs, _ := httputil.DumpResponse(this.Response, true)
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
func (this *Response) GetContentLength() string {
	return this.GetHeader("Content-Length")
}

// Cookies 获取cookie信息
func (this *Response) Cookies() (cookie []*http.Cookie) {
	if this.Response != nil {
		return this.Response.Cookies()
	}
	return
}

// CopyWith 复制数据,并监听
func (this *Response) CopyWith(w io.Writer, fn func(bs []byte)) (int, error) {
	return io.CopyWith(w, this.Body, fn)
}

// CopyWithPlan 复制数据并监听进度
func (this *Response) CopyWithPlan(w io.Writer, fn func(p *Plan)) (int, error) {
	p := &Plan{
		Current: 0,
		Total:   conv.Int64(this.GetContentLength()),
	}
	return io.CopyWith(w, this.Body, func(buf []byte) {
		p.Current += int64(len(buf))
		fn(p)
	})
}

// WriteToNewFile 写入新文件,会覆盖原文件(如果存在)
func (this *Response) WriteToNewFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return this.WriteToFile(f)
}

// WriteToFile 写入到文件,并关闭文件
func (this *Response) WriteToFile(file *os.File) error {
	_, err := this.WriteTo(file)
	return err
}

// WriteTo 写入到writer,例如文件下载,写入到文件
func (this *Response) WriteTo(writer io.Writer) (int64, error) {
	if this.Err() != nil {
		return 0, this.Err()
	}
	return io.Copy(writer, this.Response.Body)
}

// GetReadCloser 获取结果的body reader
func (this *Response) GetReadCloser() io.ReadCloser {
	return this.Response.Body
}

// GetBody 一次性读取全部字节(适用于小数据)
func (this *Response) GetBody() []byte {
	if this.body == nil && this.Response != nil {
		this.body, _ = ioutil.ReadAll(this.Response.Body)
		this.Response.Body.Close()
		//保持body可读
		this.Body = ioutil.NopCloser(bytes.NewBuffer(this.body))
	}
	return this.body
}

// GetBodyBytes 获取body内容,,返回字节
func (this *Response) GetBodyBytes() []byte {
	return this.GetBody()
}

// GetBodyString 获取body内容,返回字符串
func (this *Response) GetBodyString() string {
	return string(this.GetBodyBytes())
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

// Bind 绑定body数据,目前支持字符串,字节和json,需要指针
func (this *Response) Bind(ptr interface{}) *Response {
	body := this.GetBodyBytes()
	if ptr != nil {
		switch val := ptr.(type) {
		case *string:
			*val = string(body)
		case *[]byte:
			//val不为nil,this.body为nil可以赋值成功
			*val = body
		default:
			this.setErr(json.Unmarshal(body, ptr))
		}
	}
	return this
}

// Err 错误信息,如果错误则Response为nil
func (this *Response) Err() error {
	return this.err
}

func newResponse(req *Request, resp *http.Response, err ...error) *Response {
	r := &Response{
		Request:  req,
		Response: resp,
		err: func() error {
			if len(err) > 0 && err[0] != nil {
				return err[0]
			} else if resp != nil && resp.StatusCode/100 != 2 {
				return fmt.Errorf("状态码错误:%d", resp.StatusCode)
			}
			return nil
		}(),
	}
	if req != nil {
		r.setTryNum(req.getTry())
		r.Bind(req.bodyBind)
	}
	return r.print()
}

type Plan struct {
	Current int64 //当前数量
	Total   int64 //总数量
}
