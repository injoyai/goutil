package mux

import (
	"context"
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/middle"
	"io/fs"
	"log"
	"net/http"
	"sync"
)

type Option func(s *Server)

// WithCORS 设置跨域
func WithCORS() Option {
	return func(s *Server) {
		s.Use(func(r *Request) { middle.WithCORS(r.Writer) })
	}
}

// WithPort 设置端口
func WithPort(port ...int) Option {
	return func(s *Server) { s.SetPort(port...) }
}

// WithPrefix 设置全局前缀,注意使用
func WithPrefix(prefix string) Option {
	return func(s *Server) { s.Grouper.Prefix = prefix }
}

func New(op ...Option) *Server {
	s := &Server{
		Port:    []int{80},
		Grouper: &Grouper{Router: mux.NewRouter()},
	}
	for _, v := range op {
		v(s)
	}
	return s
}

type Server struct {
	Port []int
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

func (this *Server) Run() (err error) {

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
			err = f(port)
		}(port)
	}
	wg.Wait()
	return
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

// Static 放在最后执行,这个会占用Grouper下所有的路由
func (this *Grouper) Static(path string, dir string) *Grouper {
	path = this.Prefix + path
	s := http.StripPrefix(path, http.FileServer(http.Dir(dir)))
	this.Router.PathPrefix(path).Handler(s)
	return this
}

// StaticEmbed 放在最后执行,这个会占用Grouper下所有的路由
func (this *Grouper) StaticEmbed(path string, e embed.FS, dir string) error {
	web, err := fs.Sub(e, dir)
	if err != nil {
		return err
	}
	path = this.Prefix + path
	s := http.StripPrefix(path, http.FileServer(http.FS(web)))
	this.Router.PathPrefix(path).Handler(s)
	return nil
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
