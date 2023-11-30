package in

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/str"
	json "github.com/json-iterator/go"
	"net/http"
	"strconv"
)

type Option func(c *Client)

func WithCORS() Option {
	return func(c *Client) {
		c.AddExitOption(func(e *Exit) {
			e.SetHeaderCORS()
		})
	}
}

func WithJson() Option {
	return func(c *Client) {
		c.AddExitOption(func(e *Exit) {
			e.SetHeaderJson()
		})
	}
}

func WithDefault() Option {
	return func(c *Client) {
		c.AddExitOption(func(e *Exit) {
			e.SetCode(200)
			e.SetMark(DefaultExitMark)
			e.SetHeaderCORS()
		})
		c.SetDealResp(func(r *Resp) (httpCode int, bs []byte) {
			return http.StatusOK, r.Bytes()
		})
	}
}

func New(op ...Option) *Client {
	c := &Client{
		CodeSucc:    http.StatusOK,
		CodeFail:    http.StatusInternalServerError,
		FiledPage:   "pageNum",
		FiledSize:   "pageSize",
		DefaultSize: 10,
		PingPath:    "/ping",
	}
	for _, f := range op {
		f(c)
	}
	return c
}

type Client struct {
	ExitMark    string
	ExitOption  []ExitOption
	DealResp    func(r *Resp) (httpCode int, bs []byte)
	CodeSucc    interface{}
	CodeFail    interface{}
	FiledPage   string
	FiledSize   string
	DefaultSize uint
	PingPath    string
}

// SetDealResp 设置处理返回格式函数
func (this *Client) SetDealResp(f func(r *Resp) (httpCode int, bs []byte)) *Client {
	this.DealResp = f
	return this
}

// AddExitOption 添加退出选项
func (this *Client) AddExitOption(f ...ExitOption) *Client {
	this.ExitOption = append(this.ExitOption, f...)
	return this
}

// SetExitOption 设置退出选项
func (this *Client) SetExitOption(f ...ExitOption) *Client {
	this.ExitOption = f
	return this
}

//=================================-=================================//

// Json 返回json退出
func (this *Client) Json(httpCode int, data interface{}, count ...int64) {
	code := conv.Select(httpCode == http.StatusOK, this.CodeSucc, this.CodeFail)
	this.NewExit(httpCode, code, data, count...).SetHeaderJson().Exit()
}

// File 文件退出
func (this *Client) File(name string, bytes []byte) {
	this.NewExit(http.StatusOK, this.CodeSucc, bytes).
		SetHeader("Content-Disposition", "attachment; filename="+name).
		SetHeader("Content-Length", strconv.Itoa(len(bytes))).
		Exit()
}

// NewExit 自定义退出
func (this *Client) NewExit(httpCode int, code, data interface{}, count ...int64) *Exit {
	bs := []byte(nil)
	resp := newResp(code, data, "", count...)
	if this.DealResp != nil {
		httpCode, bs = this.DealResp(resp)
	} else {
		bs = resp.Bytes()
	}
	exit := NewExit(httpCode, bs)
	for _, v := range this.ExitOption {
		v(exit)
	}
	return exit
}

// Exit 退出
func (this *Client) Exit(httpCode int, body interface{}) {
	NewExit(httpCode, body).Exit()
}

// Succ 成功退出
func (this *Client) Succ(data interface{}, count ...int64) {
	this.Json(200, data, count...)
}

// Fail 失败退出
func (this *Client) Fail(data interface{}) {
	this.Json(500, data)
}

/***/

func (this *Client) MiddleRecover(body []byte, w http.ResponseWriter) {
	lenMark := len(this.ExitMark)
	if str.CutLeast(string(body), lenMark) == this.ExitMark {
		e := new(Exit)
		err := json.Unmarshal(body[lenMark:], &e)
		if err == nil {
			e.WriteTo(w)
			return
		}
	}
	e := this.NewExit(http.StatusInternalServerError, this.CodeFail, string(body))
	e.WriteTo(w)
}
