package in

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var DefaultClient = New(WithDefault())

func InitGo(h http.Handler) {
	DefaultClient.InitGo(h)
}

func InitGf(name ...interface{}) {
	DefaultClient.InitGf(name...)
}

func InitGin(s *gin.Engine) {
	DefaultClient.InitGin(s)
}
