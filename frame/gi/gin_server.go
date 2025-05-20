package gi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/injoyai/goutil/frame/in/v3"
	"io"
)

type (
	Context = gin.Context
)

type Server struct {
	*gin.Engine
	Port int
}

func New(port int) *Server {
	s := &Server{
		Engine: gin.Default(),
		Port:   port,
	}
	s.Use(gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recover any) {
		in.MiddleRecover(recover, c.Writer)
	}))
	return s
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

func (this *Server) Run() error {
	return this.Engine.Run(fmt.Sprintf(":%d", this.Port))
}
