package in

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/net/ghttp"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"strings"
)

//=========================Go====================

func initGo(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				body := conv.String(err)
				if strings.Contains(body, DefaultOption.ExitMark) {
					l := new(ExitModel)
					if err := json.Unmarshal([]byte(body), &l); err != nil {
						w.Write([]byte(body))
					} else {
						for i, v := range l.Header {
							w.Header().Set(i, v)
						}
						w.WriteHeader(l.Code)
						w.Write(l.Value)
					}
				} else {
					w.WriteHeader(500)
					w.Write([]byte(body))
				}
			}
		}()
		if r.URL.Path == "/ping" {
			Succ(nil)
		}
		h.ServeHTTP(w, r)
	})
}

//=========================GoFrame====================

func initGf(name ...any) *ghttp.Server {
	s := gins.Server(name...)
	s.BindStatusHandler(500, MiddleGfRecover())
	s.Group("", func(group *ghttp.RouterGroup) {
		group.ALL("/ping", func(r *ghttp.Request) { Succ(nil) })
	})
	return s
}

func MiddleGfRecover() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		body := r.Response.BufferString()
		if strings.Contains(body, DefaultOption.ExitMark) {
			l := new(ExitModel)
			if err := json.Unmarshal([]byte(body), &l); err != nil {
				r.Response.SetBuffer([]byte(body))
			} else {
				for i, v := range l.Header {
					r.Response.ResponseWriter.Header().Set(i, v)
				}
				r.Response.WriteHeader(l.Code)
				r.Response.SetBuffer(l.Value)
			}
		}
	}
}

//=========================Gin====================

func initGin(s *gin.Engine) *gin.Engine {
	s.Use(MiddleGinRecover())
	s.Any("/ping", func(c *gin.Context) { Succ(nil) })
	return s
}

func MiddleGinRecover() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recover any) {
		body := conv.String(recover)
		if strings.Contains(body, DefaultOption.ExitMark) {
			l := new(ExitModel)
			if err := json.Unmarshal([]byte(body), &l); err != nil {
				c.String(500, body)
			} else {
				for i, v := range l.Header {
					c.Header(i, v)
				}
				c.String(l.Code, string(l.Value))
			}
		} else {
			c.String(500, body)
		}
	})
}
