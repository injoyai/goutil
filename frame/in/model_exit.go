package in

import (
	"encoding/json"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/net/http"
	"strings"
)

const (
	DefaultExitMark = "未初始化!!!\n" // "EXITMARK"
)

type ExitOption func(e *Exit)

func NewExit(httpCode int, i IMarshal) *Exit {
	bs, err := i.Bytes()
	g.PanicErr(err)
	e := &Exit{
		Mark:   DefaultExitMark,
		Code:   httpCode,
		Header: http.Header{},
		Body:   bs,
	}
	e.TrySetContentType(i.ContentType()...)
	return e
}

type Exit struct {
	Mark   string      `json:"-"`      //退出标识
	Code   int         `json:"code"`   //响应状态码
	Header http.Header `json:"header"` //响应请求头
	Body   []byte      `json:"body"`   //响应内容,body
}

func (this *Exit) SetMark(mark string) *Exit {
	if len(mark) == 0 {
		mark = DefaultExitMark
	}
	this.Mark = mark
	return this
}

func (this *Exit) SetCode(code int) *Exit {
	this.Code = code
	return this
}

func (this *Exit) AddHeader(i string, v ...string) *Exit {
	this.Header[i] = append(this.Header[i], v...)
	return this
}

// SetHeader 设置请求头
func (this *Exit) SetHeader(i string, v ...string) *Exit {
	this.Header[i] = v
	return this
}

// TrySetHeader 尝试设置请求头
func (this *Exit) TrySetHeader(i string, v ...string) *Exit {
	if this.Header.Get(i) == "" {
		return this.SetHeader(i, v...)
	}
	return this
}

// SetContentType 设置请求头Content-Type
func (this *Exit) SetContentType(ct ...string) *Exit {
	return this.SetHeader(http.HeaderKeyContentType, ct...)
}

// TrySetContentType 尝试设置请求头Content-Type
func (this *Exit) TrySetContentType(ct ...string) *Exit {
	if this.Header.Get(http.HeaderKeyContentType) == "" {
		return this.SetContentType(ct...)
	}
	return this
}

// SetHeaderJson 设置请求头Content-Type
func (this *Exit) SetHeaderJson() *Exit {
	return this.SetHeader(http.HeaderKeyContentType, "application/json;charset=utf-8")
}

// TrySetHeaderJson 尝试设置请求头Content-Type
func (this *Exit) TrySetHeaderJson() *Exit {
	if this.Header.Get(http.HeaderKeyContentType) == "" {
		return this.SetHeaderJson()
	}
	return this
}

// SetHeaderCORS 设置跨域
func (this *Exit) SetHeaderCORS() *Exit {
	for k, v := range http.CORS {
		this.SetHeader(k, v...)
	}
	return this
}

// Exit 退出程序,中断执行,需要和recover配合使用
func (this *Exit) Exit() {
	bs, err := json.Marshal(this)
	g.PanicErr(err)
	panic(this.Mark + string(bs))
}

// WriteTo 写入响应
// 这里要先设置Header,再设置Code,否则个别框架设置Header可能无效(例mux)
func (this *Exit) WriteTo(w http.Writer) {
	for i, v := range this.Header {
		w.Header().Set(i, strings.Join(v, ","))
	}
	w.WriteHeader(this.Code)
	w.Write(this.Body)
}
