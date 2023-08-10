package router

import (
	"bytes"
	"net/http"
)

func newResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		Writer:     w,
		bodyBuffer: bytes.NewBuffer(nil),
		statusCode: 200,
	}
}

type ResponseWriter struct {
	Writer     http.ResponseWriter
	bodyBuffer *bytes.Buffer //响应内容
	statusCode int           //响应状态码
}

func (r *ResponseWriter) done() {
	r.Writer.WriteHeader(r.GetStatusCode())
	r.Writer.Write(r.GetBody())
}

func (r *ResponseWriter) Reset() {
	r.bodyBuffer.Reset()
}

func (r *ResponseWriter) Buffer() *bytes.Buffer {
	return r.bodyBuffer
}

func (r *ResponseWriter) GetBody() []byte {
	return r.bodyBuffer.Bytes()
}

func (r *ResponseWriter) GetBodyString() string {
	return string(r.GetBody())
}

func (r *ResponseWriter) Write(bs []byte) *ResponseWriter {
	r.bodyBuffer.Write(bs)
	return r
}

func (r *ResponseWriter) WriteBytes(bytes []byte) *ResponseWriter {
	return r.Write(bytes)
}

func (r *ResponseWriter) WriteString(s string) *ResponseWriter {
	r.bodyBuffer.WriteString(s)
	return r
}

func (r *ResponseWriter) Header() http.Header {
	return r.Writer.Header()
}

func (r *ResponseWriter) GetHeader(key string) string {
	return r.Header().Get(key)
}

func (r *ResponseWriter) GetHeaders(key string) []string {
	return r.Header().Values(key)
}

func (r *ResponseWriter) AddHeader(key, value string) {
	r.Header().Add(key, value)
}

func (r *ResponseWriter) SetHeader(key, value string) {
	r.Header().Set(key, value)
}

func (r *ResponseWriter) DelHeader(key string) {
	r.Header().Del(key)
}

func (r *ResponseWriter) GetStatusCode() int {
	return r.statusCode
}

func (r *ResponseWriter) SetStatusCode(statusCode int) *ResponseWriter {
	r.statusCode = statusCode
	return r
}
