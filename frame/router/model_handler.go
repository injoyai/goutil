package router

import (
	"reflect"
	"runtime"
)

type Handler func(r *Request)

func (h Handler) GetPath() string {
	if h == nil {
		return "nil"
	}
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

type MiddleFunc []Handler

func (m *MiddleFunc) Add(handler Handler) {
	*m = append(*m, handler)
}

func (m MiddleFunc) Do(r *Request) {
	for _, fn := range m {
		fn(r)
	}
}

// DefaultHandle400 Bad Request
func DefaultHandle400() Handler {
	return func(r *Request) {
		r.SetStatusCode(400)
		r.ClearBody()
		r.WriteString("Bad Request")
	}
}

// DefaultHandle401 Unauthorized
func DefaultHandle401() Handler {
	return func(r *Request) {
		r.SetStatusCode(401)
		r.ClearBody()
		r.WriteString("Unauthorized")
	}
}

// DefaultHandle403 Forbidden
func DefaultHandle403() Handler {
	return func(r *Request) {
		r.SetStatusCode(403)
		r.ClearBody()
		r.WriteString("Forbidden")
	}
}

// DefaultHandle404 Page Not Find
func DefaultHandle404() Handler {
	return func(r *Request) {
		r.SetStatusCode(404)
		r.ClearBody()
		r.WriteString("Page Not Find")
	}
}

// DefaultHandle405 Method Not Allowed
func DefaultHandle405() Handler {
	return func(r *Request) {
		r.SetStatusCode(405)
		r.ClearBody()
		r.WriteString("Method Not Allowed")
	}
}

// DefaultHandle500 Internal Server Error
func DefaultHandle500() Handler {
	return func(r *Request) {}
}
