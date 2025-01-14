package gf

import (
	"fmt"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/net/ghttp"
	"github.com/injoyai/goutil/frame/gf/swagger"
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/i/html"
	"github.com/injoyai/goutil/net/ip"
	"net/http"
	"time"
)

type Server struct {
	*ghttp.Server
	Port        []int //端口
	enablePprof bool  //
}

/*

//cpu检测
go tool pprof http://localhost:6060/pprof/profile?seconds=20

//内存检测
go tool pprof http://localhost:6060/pprof/heap

*/

// New 快速开始默认配置
// @port,端口
func New(name ...interface{}) *Server {
	s := gins.Server(name...)
	s.SetClientMaxBodySize(8 << 20) //设置body最大数据,8m
	s.SetAccessLogEnabled(false)    //请求日志
	s.SetErrorLogEnabled(false)     //错误日志
	s.BindStatusHandler(http.StatusNotFound, func(r *ghttp.Request) {
		r.Response.ClearBuffer()
		r.Response.WriteExit(html.PageNotFindRobot)
	})
	s.BindStatusHandler(http.StatusInternalServerError, func(r *ghttp.Request) {
		body := r.Response.Buffer()
		r.Response.ClearBuffer()
		in.MiddleRecover(body, r.Response.Writer)
	})
	return &Server{Server: s}
}

func (this *Server) SetPort(port ...int) *Server {
	this.Port = port
	this.Server.SetPort(port...)
	return this
}

func (this *Server) EnablePProf(pattern ...string) *Server {
	this.enablePprof = true
	this.Server.EnablePProf(pattern...)
	return this
}

func (this *Server) SetSwagger(prefix, path string) *Server {
	this.Plugin(&swagger.Swagger{
		Prefix: prefix,
		Path:   path,
	})
	return this
}

func (this *Server) SetShowRoute(b ...bool) *Server {
	this.Server.SetDumpRouterMap(!(len(b) > 0 && !b[0]))
	return this
}

func (this *Server) UseCORS() *Server {
	this.Use(func(r *ghttp.Request) {
		r.Response.CORSDefault()
		r.Middleware.Next()
	})
	return this
}

func (this *Server) ALL(s string, fn func(r *ghttp.Request)) *Server {
	this.Server.Group("", func(g *ghttp.RouterGroup) { g.ALL(s, fn) })
	return this
}

func (this *Server) GET(s string, fn func(r *ghttp.Request)) *Server {
	this.Server.Group("", func(g *ghttp.RouterGroup) { g.GET(s, fn) })
	return this
}

func (this *Server) POST(s string, fn func(r *ghttp.Request)) *Server {
	this.Server.Group("", func(g *ghttp.RouterGroup) { g.POST(s, fn) })
	return this
}

func (this *Server) PUT(s string, fn func(r *ghttp.Request)) *Server {
	this.Server.Group("", func(g *ghttp.RouterGroup) { g.PUT(s, fn) })
	return this
}

func (this *Server) DELETE(s string, fn func(r *ghttp.Request)) *Server {
	this.Server.Group("", func(g *ghttp.RouterGroup) { g.DELETE(s, fn) })
	return this
}

func (this *Server) Run() {
	go func() {
		<-time.After(time.Millisecond * 100)
		for _, port := range this.Port {
			ipv4 := ip.GetLocal()
			fmt.Printf("打开接口文档: 点击 http://%s:%d/swagger 	\n", ipv4, port)
			if this.enablePprof {
				fmt.Printf("打开性能剖析: 点击 http://%s:%d/pprof 	\n", ipv4, port)
			}
			break
		}
	}()
	this.Server.Run()
}
