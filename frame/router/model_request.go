package router

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func newRequest(w http.ResponseWriter, r *http.Request) *Request {
	return &Request{
		Request:        r,
		requestBody:    nil,
		ResponseWriter: newResponseWriter(w),
		Middle:         &Middle{},
	}
}

type Request struct {
	*http.Request           //请求
	requestBody     []byte  //请求内容
	*ResponseWriter         //响应
	Middle          *Middle //中间件
}

func (r *Request) RequestBody() io.ReadCloser {
	return r.Request.Body
}

func (r *Request) GetBody() []byte {
	if r.requestBody == nil {
		r.requestBody, _ = ioutil.ReadAll(r.Request.Body)
		r.Request.Body.Close()
	}
	return r.requestBody
}

func (r *Request) GetBodyString() string {
	return string(r.GetBody())
}

func (r *Request) GetQueryString(key string, def ...string) string {
	if val, ok := r.Request.URL.Query()[key]; ok {
		return val[0]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (r *Request) GetQuery() url.Values {
	return r.Request.URL.Query()
}

func (r *Request) UnmarshalJson(v any) error {
	return json.Unmarshal(r.GetBody(), v)
}

func (r *Request) WriteErr(err error) *Request {
	if err != nil {
		r.ClearBody()
		r.WriteString(err.Error())
		r.SetStatusCode(500)
	}
	return r
}

func (r *Request) SetStatusCode(code int) *Request {
	r.ResponseWriter.SetStatusCode(code)
	return r
}

func (r *Request) ResetBody() *Request {
	r.ResponseWriter.Reset()
	return r
}

func (r *Request) ClearBody() *Request {
	return r.ResetBody()
}

func (r *Request) Write(bs []byte) *Request {
	r.ResponseWriter.Write(bs)
	return r
}

func (r *Request) WriteBytes(bs []byte) *Request {
	return r.Write(bs)
}

func (r *Request) WriteBytesExit(bs []byte) {
	r.Write(bs).Exit()
}

func (r *Request) WriteString(s string) *Request {
	r.ResponseWriter.WriteString(s)
	return r
}

func (r *Request) WriteStringExit(s string) {
	r.WriteString(s).Exit()
}

func (r *Request) WriteJson(i any) *Request {
	bs, _ := json.Marshal(i)
	r.SetHeader("Content-Type", "application/json;charset=utf-8")
	r.WriteBytes(bs)
	return r
}

func (r *Request) WriteJsonExit(i any) {
	r.WriteJson(i).Exit()
}

func (r *Request) SetHeader(key, val string) {
	r.ResponseWriter.SetHeader(key, val)
}

func (r *Request) AddHeader(key, val string) {
	r.ResponseWriter.AddHeader(key, val)
}

func (r *Request) DelHeader(key string) {
	r.ResponseWriter.DelHeader(key)
}

func (r *Request) GetHeader(key string) string {
	return r.Request.Header.Get(key)
}

func (r *Request) GetHeaders(key string) []string {
	return r.Request.Header.Values(key)
}

func (r *Request) done() {
	r.ResponseWriter.done()
}

func (r *Request) Exit() {
	panic(MarkExit)
}
