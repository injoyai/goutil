package in

import (
	"encoding/json"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"strconv"
	"strings"
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
		c.SetSuccFailCode(http.StatusOK, http.StatusInternalServerError)
	}
}

// WithQJ 设置启进响应数据格式
func WithQJ() Option {
	return func(c *Client) {
		c.AddExitOption(func(e *Exit) {
			e.SetMark(DefaultExitMark)
			e.SetHeaderCORS()
		})
		c.SetSuccFailCode("SUCCESS", "FAIL")
	}
}

func New(op ...Option) *Client {
	c := &Client{
		ExitMark:    DefaultExitMark,
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
	FiledPage   string
	FiledSize   string
	DefaultSize uint
	PingPath    string
	SuccFail    func(c *Client, succ bool, data interface{}, count ...int64)
}

// SetSuccFailCode 设置响应成功失败
func (this *Client) SetSuccFailCode(succ, fail interface{}) *Client {
	return this.SetSuccFail(func(c *Client, ok bool, data interface{}, count ...int64) {
		if ok {
			c.Json(http.StatusOK, NewRespMap(succ, data, count...))
		} else {
			c.Json(http.StatusOK, NewRespMap(fail, data, count...))
		}
	})
}

func (this *Client) SetSuccFail(f func(c *Client, succ bool, data interface{}, count ...int64)) *Client {
	this.SuccFail = f
	return this
}

func (this *Client) SetPingPath(path string) *Client {
	this.PingPath = path
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

//=================================Response=================================//

func (this *Client) Redirect(httpCode int, url string) {
	this.NewExit(httpCode, &TEXT{}).SetHeader("Location", url).Exit()
}

// File 文件退出
func (this *Client) File(name string, bytes []byte) {
	this.NewExit(http.StatusOK, &TEXT{Data: bytes}).
		SetHeader("Content-Disposition", "attachment; filename="+name).
		SetHeader("Content-Length", strconv.Itoa(len(bytes))).
		Exit()
}

// Json 返回json退出
func (this *Client) Json(httpCode int, data interface{}) {
	this.NewExit(httpCode, &JSON{Data: data}).Exit()
}

func (this *Client) Html(httpCode int, data interface{}) {
	this.NewExit(httpCode, &HTML{Data: data}).Exit()
}

func (this *Client) Text(httpCode int, data interface{}) {
	this.NewExit(httpCode, &TEXT{Data: data}).Exit()
}

// NewExit 自定义退出
func (this *Client) NewExit(httpCode int, i IMarshal) *Exit {
	exit := NewExit(httpCode, i)
	for _, v := range this.ExitOption {
		v(exit)
	}
	return exit
}

//=================================Other=================================//

// Succ 成功退出,自定义
func (this *Client) Succ(data interface{}, count ...int64) {
	if this.SuccFail != nil {
		this.SuccFail(this, true, data, count...)
		return
	}
	this.Json(http.StatusOK, NewRespMap(http.StatusOK, data, count...))
}

// Fail 失败退出,自定义
func (this *Client) Fail(data interface{}) {
	if this.SuccFail != nil {
		this.SuccFail(this, false, data)
		return
	}
	this.Json(http.StatusOK, NewRespMap(http.StatusInternalServerError, data))
}

func (this *Client) CopyFile(w http.ResponseWriter, name string, r io.Reader) {
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	this.Copy(w, r)
}

func (this *Client) Copy(w http.ResponseWriter, r io.Reader) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, r)
}

func (this *Client) Proxy(w http.ResponseWriter, r *http.Request, uri string) {
	req, err := http.NewRequest(r.Method, uri, r.Body)
	if err != nil {
		this.Fail(err)
		return
	}
	req.Header = r.Header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		this.Fail(err)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

//=================================Middle=================================//

// Recover 初始化原生
func (this *Client) Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				this.MiddleRecover(err, w)
			}
		}()
		if len(this.PingPath) > 0 && r.URL.Path == this.PingPath {
			this.Succ(nil)
		}
		h.ServeHTTP(w, r)
	})
}

func (this *Client) MiddleRecover(err interface{}, w http.ResponseWriter) {
	bs := []byte(conv.String(err))
	if strings.HasPrefix(string(bs), this.ExitMark) {
		e := new(Exit)
		err := json.Unmarshal(bs[len(this.ExitMark):], &e)
		if err == nil {
			e.WriteTo(w)
			return
		}
	}
	this.NewExit(http.StatusInternalServerError, &TEXT{Data: err}).WriteTo(w)
}

func (this *Client) GetPageNum(r interface{}) int {
	return Get(r, this.FiledPage, 1).Int() - 1
}

func (this *Client) GetPageSize(r interface{}) int {
	return Get(r, this.FiledSize, this.DefaultSize).Int()
}
