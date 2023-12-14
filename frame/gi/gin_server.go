package gi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/injoyai/goutil/frame/in/v2"
	"github.com/injoyai/logs"
)

type (
	Context = gin.Context
)

type Server struct {
	*gin.Engine
	Port int
}

func New(port int) *Server {
	data := &Server{
		Engine: in.InitGin(gin.Default()),
		Port:   port,
	}
	return data
}

func (this *Server) ALL(s string, fn gin.HandlerFunc) *Server {
	this.Engine.Any(s, fn)
	return this
}

func (this *Server) GET(s string, fn gin.HandlerFunc) *Server {
	this.Engine.GET(s, fn)
	return this
}

func (this *Server) POST(s string, fn gin.HandlerFunc) *Server {
	this.Engine.POST(s, fn)
	return this
}

func (this *Server) PUT(s string, fn gin.HandlerFunc) *Server {
	this.Engine.PUT(s, fn)
	return this
}

func (this *Server) DELETE(s string, fn gin.HandlerFunc) *Server {
	this.Engine.DELETE(s, fn)
	return this
}

func (this *Server) Run(ports ...int) {
	addr := []string(nil)
	for _, v := range ports {
		addr = append(addr, fmt.Sprintf(":%d", v))
	}
	logs.PrintErr(this.Engine.Run(addr...))
}
