package in

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

var DefaultClient = New(WithDefault())

func InitGo(h http.Handler) http.Handler {
	return DefaultClient.InitGo(h)
}

func InitGf(name ...interface{}) *ghttp.Server {
	return DefaultClient.InitGf(name...)
}

func InitGin(s *gin.Engine) *gin.Engine {
	return DefaultClient.InitGin(s)
}
