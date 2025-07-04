package in

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"strings"
	"unsafe"
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
		c.SetHandlerWithCode(
			http.StatusOK,
			http.StatusInternalServerError,
			http.StatusUnauthorized,
			http.StatusForbidden,
		)
	}
}

// WithQJ 设置启进响应数据格式
func WithQJ() Option {
	return func(c *Client) {
		c.SetExitOption(func(e *Exit) {
			e.SetHeaderCORS()
		})
		c.SetHandlerWithCode("SUCCESS", "FAIL", "FAIL", "FAIL")
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
	c.SetHandlerWithCode(
		http.StatusOK,
		http.StatusInternalServerError,
		http.StatusUnauthorized,
		http.StatusForbidden,
	)
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
	BindMap             map[string]http.HandlerFunc    //自定义接口绑定
	HandlerSucc         func(data any, count ...int64) //成功
	HandlerFail         func(data any)                 //失败
	HandlerUnauthorized func()                         //验证失败
	HandlerForbidden    func()                         //权限不足
}

// SetHandlerWithCode 设置响应成功失败等
func (this *Client) SetHandlerWithCode(succ, fail, unauthorized, forbidden any) *Client {
	this.HandlerSucc = this.NewSuccWithCode(succ)
	this.HandlerFail = this.NewFailWithCode(fail)
	this.HandlerUnauthorized = this.NewUnauthorizedWithCode(unauthorized)
	this.HandlerForbidden = this.NewForbiddenWithCode(forbidden)
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
func (this *Client) Json(httpCode int, data any) {
	this.NewExit(httpCode, &JSON{Data: data}).Exit()
}

func (this *Client) Html(httpCode int, data any) {
	this.NewExit(httpCode, &HTML{Data: data}).Exit()
}

func (this *Client) Text(httpCode int, data any) {
	this.NewExit(httpCode, &TEXT{Data: data}).Exit()
}

func (this *Client) File(name string, size int64, r io.ReadCloser) {
	this.NewExit(http.StatusOK, &FILE{
		Name:       name,
		Size:       size,
		ReadCloser: r}).Exit()
}

func (this *Client) Reader(httpCode int, r io.ReadCloser) {
	this.NewExit(httpCode, &READER{ReadCloser: r}).Exit()
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
func (this *Client) Succ(data any, count ...int64) {
	if this.HandlerSucc == nil {
		this.HandlerSucc = this.NewSuccWithCode(http.StatusOK)
	}
	this.HandlerSucc(data, count...)
}

// Fail 失败退出,自定义
func (this *Client) Fail(msg any) {
	if this.HandlerFail == nil {
		this.HandlerFail = this.NewFailWithCode(http.StatusInternalServerError)
	}
	this.HandlerFail(msg)
}

func (this *Client) Unauthorized() {
	if this.HandlerUnauthorized == nil {
		this.HandlerUnauthorized = this.NewUnauthorizedWithCode(http.StatusUnauthorized)
	}
	this.HandlerUnauthorized()
}

func (this *Client) Forbidden() {
	if this.HandlerForbidden == nil {
		this.HandlerForbidden = this.NewForbiddenWithCode(http.StatusForbidden)
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

// MiddleRecover 例gf等web框架只需要这一半即可,但是Bind会失效
func (this *Client) MiddleRecover(e any, w http.ResponseWriter) {
	switch w2 := e.(type) {
	case *Exit:
		w2.WriteTo(w)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		s := conv.String(e)
		w.Write(*(*[]byte)(unsafe.Pointer(&s)))
	}
}

//=================================Other=================================//

func (this *Client) GetPageNum(r *http.Request) int {
	if v, ok := r.URL.Query()[this.FiledPage]; ok && len(v) > 0 {
		return conv.Int(v[0]) - 1
	}
	return 0
}

func (this *Client) GetPageSize(r *http.Request) int {
	if v, ok := r.URL.Query()[this.FiledSize]; ok && len(v) > 0 {
		return conv.Int(v)
	}
	return this.DefaultSize
}
