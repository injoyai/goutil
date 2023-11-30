package in

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/net/ghttp"
	"github.com/injoyai/conv"
	"io"
	"net/http"
)

// InitGo 初始化原生
func (this *Client) InitGo(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				this.MiddleRecover([]byte(conv.String(err)), w)
			}
		}()
		if r.URL.Path == this.PingPath {
			this.Succ(nil)
		}
		h.ServeHTTP(w, r)
	})
}

// InitGf 初始化GoFrame
func (this *Client) InitGf(name ...interface{}) *ghttp.Server {
	s := gins.Server(name...)
	s.BindStatusHandler(http.StatusInternalServerError, func(r *ghttp.Request) {
		body := r.Response.Buffer()
		this.MiddleRecover(body, r.Response.Writer)
	})
	s.Group("", func(group *ghttp.RouterGroup) {
		group.ALL(this.PingPath, func(r *ghttp.Request) { this.Succ(nil) })
	})
	return s
}

// InitGin 初始化Gin
func (this *Client) InitGin(s *gin.Engine) *gin.Engine {
	s.Use(gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recover interface{}) {
		body := []byte(conv.String(recover))
		this.MiddleRecover(body, c.Writer)
	}))
	s.Any(this.PingPath, func(c *gin.Context) { this.Succ(nil) })
	return s
}
