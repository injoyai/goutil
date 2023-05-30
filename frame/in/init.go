package in

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

var (
	DefaultOption     = defaultOption()
	DefaultFunc       = defaultFunc()
	ErrInvalidRequest = errors.New("invalid request")
)

func init() {

	//原生需要手动引用Init()

	//gin需要手动引用

	//初始化gf框架
	initGf()

}

// InitGf 初始化gf框架
//除了单例模式的其余需要手动初始化
func InitGf(name ...interface{}) *ghttp.Server {
	s := gins.Server(name...)
	if s.GetName() != "default" {
		initGf(name...)
	}
	return s
}

/*
	例:
	http.ListenAndServe(":8001", Handler)
	TO
	http.ListenAndServe(":8001", in.InitGo(Handler))
*/
// InitGo 原始框架需要引用此方法
func InitGo(h http.Handler) http.Handler {
	return initGo(h)
}

// ListenAndServe 原始服务初始化
func ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, InitGo(handler))
}

// ListenAndServeTLS 原始服务初始化
func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, InitGo(handler))
}

// InitGin gin框架需要手动引用
func InitGin(r *gin.Engine) *gin.Engine {
	return initGin(r)
}
