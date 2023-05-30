package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/injoyai/goutil/frame/in"
	"github.com/injoyai/logs"
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
	this.Any(s, fn)
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

func (this *Server) Run(port ...int) {
	if len(port) > 0 {
		this.Port = port[0]
	}
	logs.PrintErr(this.Engine.Run(fmt.Sprintf(":%d", this.Port)))
}
