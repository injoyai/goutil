package in

import "net/http"

// Recover 初始化原生
func (this *Client) Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				this.MiddleRecover(err, w)
			}
		}()
		if len(this.PingPath) > 0 && r.URL.Path == this.PingPath {
			this.Succ(nil)
		}
		h.ServeHTTP(w, r)
	})
}

//// InitGf 初始化GoFrame
//func (this *Client) RecoverGf(name ...any) *ghttp.Server {
//	s := gins.Server(name...)
//	s.BindStatusHandler(http.StatusInternalServerError, func(r *ghttp.Request) {
//		body := r.Response.Buffer()
//		r.Response.ClearBuffer()
//		this.MiddleRecover(body, r.Response.Writer)
//	})
//	if len(this.PingPath) > 0 {
//		s.Group("", func(group *ghttp.RouterGroup) {
//			group.ALL(this.PingPath, func(r *ghttp.Request) { this.Succ(nil) })
//		})
//	}
//	return s
//}

//// InitGin 初始化Gin
//func (this *Client) RecoverGin(s *gin.Engine) *gin.Engine {
//	s.Use(gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recover any) {
//		this.MiddleRecover(recover, c.Writer)
//	}))
//	if len(this.PingPath) > 0 {
//		s.Any(this.PingPath, func(c *gin.Context) { this.Succ(nil) })
//	}
//	return s
//}
