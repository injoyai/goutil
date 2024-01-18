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

func SetSuccFailCode(succ, fail interface{}) *Client {
	return DefaultClient.SetSuccFailCode(succ, fail)
}

func SetSuccFail(f func(c *Client, succ bool, data interface{}, count ...int64)) *Client {
	return DefaultClient.SetSuccFail(f)
}
