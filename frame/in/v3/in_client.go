package in

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"strings"
)

type Option func(c *Client)

// WithCORS 设置响应CORS头部
func WithCORS() Option {
	return func(c *Client) {
		c.SetExitOption(func(e *Exit) {
			e.SetHeaderCORS()
		})
	}
}

// WithJson 设置响应json头部
func WithJson() Option {
	return func(c *Client) {
		c.SetExitOption(func(e *Exit) {
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
		c.SetExitOption(func(e *Exit) {
			e.SetHeaderCORS()
		})
		c.SetStatusCode(
			http.StatusOK,
			http.StatusInternalServerError,
			http.StatusUnauthorized,
			http.StatusForbidden)
	}
}

// WithQJ 设置启进响应数据格式
func WithQJ() Option {
	return func(c *Client) {
		c.SetExitOption(func(e *Exit) {
			e.SetHeaderCORS()
		})
		c.SetStatusCode("SUCCESS", "FAIL", "FAIL", "FAIL")
		c.HandlerUnauthorized = func() { c.Text(http.StatusUnauthorized, "验证失败") }
		c.HandlerForbidden = func() { c.Text(http.StatusForbidden, "没有权限") }
	}
}

func New(op ...Option) *Client {
	c := &Client{
		Safe:        maps.NewSafe(),
		FiledPage:   "pageNum",
		FiledSize:   "pageSize",
		DefaultSize: 10,
		BindMap:     map[string]http.HandlerFunc{},
	}
	c.Bind("/ping", func(w http.ResponseWriter, r *http.Request) { c.Succ(nil) })
	c.SetStatusCode(
		http.StatusOK,
		http.StatusInternalServerError,
		http.StatusUnauthorized,
		http.StatusForbidden)
	for _, f := range op {
		f(c)
	}
	return c
}

type Client struct {
	*maps.Safe
	ExitOption          []ExitOption
	FiledPage           string
	FiledSize           string
	DefaultSize         int
	BindMap             map[string]http.HandlerFunc            //自定义接口绑定
	HandlerSucc         func(data interface{}, count ...int64) //成功
	HandlerFail         func(data interface{})                 //失败
	HandlerUnauthorized func()                                 //验证失败
	HandlerForbidden    func()                                 //权限不足
}

// SetStatusCode 设置响应成功失败
func (this *Client) SetStatusCode(succ, fail, unauthorized, forbidden interface{}) *Client {
	this.HandlerSucc = func(data interface{}, count ...int64) {
		if len(count) > 0 {
			this.Json(http.StatusOK, &ResponseCount{
				Code:    succ,
				Data:    data,
				Message: "成功",
				Count:   count[0],
			})
			return
		}
		this.Json(http.StatusOK, &Response{
			Code:    succ,
			Data:    data,
			Message: "成功",
		})
	}
	this.HandlerFail = func(msg interface{}) {
		this.Json(http.StatusOK, &Response{
			Code:    fail,
			Message: conv.String(msg),
		})
	}
	this.HandlerUnauthorized = func() {
		this.Json(http.StatusOK, &Response{
			Code:    unauthorized,
			Message: "验证失败",
		})
	}
	this.HandlerForbidden = func() {
		this.Json(http.StatusOK, &Response{
			Code:    forbidden,
			Message: "权限不足",
		})
	}
	return this
}

func (this *Client) Bind(path string, handler http.HandlerFunc) *Client {
	this.BindMap[path] = handler
	return this
}

// SetExitOption 设置退出选项
func (this *Client) SetExitOption(f ...ExitOption) *Client {
	this.ExitOption = append(this.ExitOption, f...)
	return this
}

//=================================Response=================================//

func (this *Client) Redirect(httpCode int, url string) {
	this.NewExit(httpCode, &TEXT{}).SetHeader("Location", url).Exit()
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

func (this *Client) File(name string, size int64, r io.ReadCloser) {
	this.NewExit(http.StatusOK, &FILE{
		Name:       name,
		Size:       size,
		ReadCloser: r}).Exit()
}

// NewExit 自定义退出
func (this *Client) NewExit(httpCode int, i IMarshal) *Exit {
	return NewExit(httpCode, i, this.ExitOption...)
}

// Exit 直接退出,设置的跨域啥的应该是无效的
func (this *Client) Exit() {
	NewExit(-1, nil).Exit()
}

//=================================Other=================================//

// Succ 成功退出,自定义
func (this *Client) Succ(data interface{}, count ...int64) {
	if this.HandlerSucc == nil {
		if len(count) > 0 {
			this.Json(http.StatusOK, &ResponseCount{
				Code:    http.StatusOK,
				Data:    data,
				Message: "成功",
				Count:   count[0],
			})
			return
		}
		this.Json(http.StatusOK, &Response{
			Code:    http.StatusOK,
			Data:    data,
			Message: "成功",
		})
	}
	this.HandlerSucc(data, count...)
}

// Fail 失败退出,自定义
func (this *Client) Fail(msg interface{}) {
	if this.HandlerFail == nil {
		this.Json(http.StatusOK, &Response{
			Code:    http.StatusInternalServerError,
			Message: conv.String(msg),
		})
	}
	this.HandlerFail(msg)
}

func (this *Client) Unauthorized() {
	if this.HandlerUnauthorized == nil {
		this.Json(http.StatusOK, &Response{
			Code:    http.StatusUnauthorized,
			Message: "验证失败",
		})
	}
	this.HandlerUnauthorized()
}

func (this *Client) Forbidden() {
	if this.HandlerForbidden == nil {
		this.Json(http.StatusOK, &Response{
			Code:    http.StatusForbidden,
			Message: "没有权限",
		})
	}
	this.HandlerForbidden()
}

func (this *Client) Proxy(w http.ResponseWriter, r *http.Request, uri string) {
	defer r.Body.Close()
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

func (this *Client) Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				this.MiddleRecover(e, w)
			}
		}()
		if handler, ok := this.BindMap[r.URL.Path]; ok && handler != nil {
			handler(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// MiddleRecover 例gf等web框架只需要这一半即可
func (this *Client) MiddleRecover(e interface{}, w http.ResponseWriter) {
	switch w2 := e.(type) {
	case *Exit:
		w2.WriteTo(w)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(conv.String(e)))
	}
}

//=================================Other=================================//

func (this *Client) GetPageNum(r *http.Request) int {
	if v, ok := r.URL.Query()[this.FiledPage]; ok {
		return conv.Int(v) - 1
	}
	return 0
}

func (this *Client) GetPageSize(r *http.Request) int {
	if v, ok := r.URL.Query()[this.FiledSize]; ok {
		return conv.Int(v)
	}
	return this.DefaultSize
}
