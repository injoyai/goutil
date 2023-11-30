package in

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	json "github.com/json-iterator/go"
	"strings"
)

const (
	DefaultExitMark = "EXIT_MARK"
)

type ExitOption func(e *Exit)

func NewExit(code int, body interface{}) *Exit {
	return &Exit{
		Mark:   DefaultExitMark,
		Code:   code,
		Header: http.Header{},
		Body:   conv.Bytes(body),
	}
}

type Exit struct {
	Mark   string      //退出标识
	Code   int         //响应状态码
	Header http.Header //响应请求头
	Body   []byte      //响应内容,body
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

// SetHeaderJson 设置请求头Content-Type
func (this *Exit) SetHeaderJson() *Exit {
	return this.SetHeader(http.HeaderKeyContentType, "application/json;charset=utf-8")
}

// SetHeaderCORS 设置跨域
func (this *Exit) SetHeaderCORS() *Exit {
	for k, v := range http.CORS {
		this.SetHeader(k, v...)
	}
	return this
}

func (this *Exit) String() string {
	bs, _ := json.Marshal(this)
	return string(bs)
}

// Exit 退出程序,中断执行,需要和recover配合使用
func (this *Exit) Exit() {
	panic(this.String())
}

func (this *Exit) WriteTo(w http.Writer) {
	w.WriteHeader(this.Code)
	for i, v := range this.Header {
		w.Header().Set(i, strings.Join(v, ","))
	}
	w.Write(this.Body)
}
