package mux

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	in "github.com/injoyai/goutil/frame/in/mini"
	"log"
	"net/http"
	"sync"
)

func New() *Server {
	s := &Server{
		Port:    []int{80},
		Server:  &http.Server{},
		Grouper: &Grouper{Router: mux.NewRouter()},
	}
	return s
}

type Server struct {
	Port   []int
	Server *http.Server

	*Grouper
}

func (this *Server) SetPort(port ...int) *Server {
	this.Port = port
	return this
}

func (this *Server) Use(f ...func(r *Request)) {
	this.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req := NewRequest(w, r)
			for _, v := range f {
				v(req)
			}
			r = r.WithContext(context.WithValue(r.Context(), "_cache", req.cache))
			next.ServeHTTP(w, r)
		})
	})
}

func (this *Server) Run() error {

	f := func(port int) (err error) {
		defer func() {
			log.Println("HTTP服务结束监听:", err)
		}()
		s := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: in.Recover(this.Router),
		}
		log.Println("HTTP服务开启监听", s.Addr)
		return s.ListenAndServe()
	}

	if len(this.Port) == 0 {
		this.Port = []int{80}
	}

	wg := sync.WaitGroup{}
	for _, port := range this.Port {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			f(port)
		}(port)
	}
	wg.Wait()
	return nil
}

type Grouper struct {
	Router *mux.Router        //路由实例
	Prefix string             //路由前缀
	middle []func(r *Request) //中间件
}

func (this *Grouper) Middle(f ...func(r *Request)) {
	this.middle = append(this.middle, f...)
}

func (this *Grouper) IgnoreMiddle(handler func(g *Grouper)) *Grouper {
	g := *this
	x := &g
	if handler != nil {
		handler(x)
	}
	return x
}

func (this *Grouper) Group(path string, handler func(g *Grouper)) *Grouper {
	prefix := this.Prefix + path
	g := &Grouper{
		Router: this.Router,
		Prefix: prefix,
		middle: this.middle,
	}
	if handler != nil {
		handler(g)
	}
	return g
}

func (this *Grouper) ALL(path string, handler func(r *Request)) *Grouper {
	path = this.Prefix + path
	middle := this.middle
	this.Router.Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodPatch, http.MethodConnect, http.MethodOptions, http.MethodTrace).
		Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(w, r)
		for _, f := range middle {
			f(req)
		}
		handler(req)
	})
	return this
}

func (this *Grouper) GET(path string, handler func(r *Request)) *Grouper {
	path = this.Prefix + path
	middle := this.middle
	this.Router.Methods(http.MethodGet).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(w, r)
		for _, f := range middle {
			f(req)
		}
		handler(req)
	})
	return this
}

func (this *Grouper) POST(path string, handler func(r *Request)) *Grouper {
	path = this.Prefix + path
	middle := this.middle
	this.Router.Methods(http.MethodPost).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(w, r)
		for _, f := range middle {
			f(req)
		}
		handler(req)
	})
	return this
}

func (this *Grouper) PUT(path string, handler func(r *Request)) *Grouper {
	path = this.Prefix + path
	middle := this.middle
	this.Router.Methods(http.MethodPut).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(w, r)
		for _, f := range middle {
			f(req)
		}
		handler(req)
	})
	return this
}

func (this *Grouper) DELETE(path string, handler func(r *Request)) *Grouper {
	path = this.Prefix + path
	middle := this.middle
	this.Router.Methods(http.MethodDelete).Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(w, r)
		for _, f := range middle {
			f(req)
		}
		handler(req)
	})
	return this
}
