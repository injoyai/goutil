package in

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/str"
	json "github.com/json-iterator/go"
	"net/http"
	"strconv"
)

type Option func(c *Client)

// WithCORS 设置响应CORS头部
func WithCORS() Option {
	return func(c *Client) {
		c.AddExitOption(func(e *Exit) {
			e.SetHeaderCORS()
		})
	}
}

// WithJson 设置响应json头部
func WithJson() Option {
	return func(c *Client) {
		c.AddExitOption(func(e *Exit) {
			e.SetHeaderJson()
		})
	}
}

// WithDefault 默认
func WithDefault() Option {
	return WithQL()
}

// WithQL 设置钱浪响应数据格式
func WithQL() Option {
	return func(c *Client) {
		c.AddExitOption(func(e *Exit) {
			e.SetMark(DefaultExitMark)
			e.SetHeaderCORS()
		})
		c.SetDealResp(func(r *Resp) (httpCode int, bs []byte) {
			return r.Default()
		})
	}
}

// WithQJ 设置启进响应数据格式
func WithQJ() Option {
	return func(c *Client) {
		c.AddExitOption(func(e *Exit) {
			e.SetMark(DefaultExitMark)
			e.SetHeaderCORS()
		})
		c.SetDealResp(func(r *Resp) (httpCode int, bs []byte) {
			switch r.Code {
			case http.StatusOK:
				r.Code = "SUCCESS"
			case http.StatusInternalServerError:
				r.Code = "FAIL"
			default:
				r.Code = "FAIL"
				return r.Code.(int), r.Bytes()
			}
			return r.Default()
		})
	}
}

func New(op ...Option) *Client {
	c := &Client{
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
	FiledPage   string
	FiledSize   string
	DefaultSize uint
	PingPath    string
}

// SetCode 设置响应成功失败
func (this *Client) SetCode(succ, fail interface{}) *Client {
	return this.SetDealResp(func(r *Resp) (httpCode int, bs []byte) {
		switch r.Code {
		case http.StatusOK:
			r.Code = succ
		case http.StatusInternalServerError:
			r.Code = fail
		}
		return r.Default()
	})
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
	this.NewExit(httpCode, data, count...).SetHeaderJson().Exit()
}

// File 文件退出
func (this *Client) File(name string, bytes []byte) {
	this.NewExit(http.StatusOK, bytes).
		SetHeader("Content-Disposition", "attachment; filename="+name).
		SetHeader("Content-Length", strconv.Itoa(len(bytes))).
		Exit()
}

// NewExit 自定义退出
func (this *Client) NewExit(httpCode int, data interface{}, count ...int64) *Exit {
	bs := []byte(nil)
	resp := newResp(httpCode, data, count...)
	if this.DealResp != nil {
		httpCode, bs = this.DealResp(resp)
	} else {
		bs, _ = json.Marshal(data)
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
	this.Json(http.StatusOK, data, count...)
}

// Fail 失败退出
func (this *Client) Fail(data interface{}) {
	this.Json(http.StatusInternalServerError, data)
}

/***/

func (this *Client) MiddleRecover(err interface{}, w http.ResponseWriter) {
	lenMark := len(this.ExitMark)
	bs := []byte(conv.String(err))
	if str.CutLeast(string(bs), lenMark) == this.ExitMark {
		e := new(Exit)
		err := json.Unmarshal(bs[lenMark:], &e)
		if err == nil {
			e.WriteTo(w)
			return
		}
	}
	e := this.NewExit(http.StatusInternalServerError, err)
	e.SetHeaderJson().WriteTo(w)
}

func (this *Client) GetPageSize(r interface{}) (int, int) {
	return Get(r, this.FiledPage, 1).Int() - 1,
		Get(r, this.FiledSize, this.DefaultSize).Int()
}
