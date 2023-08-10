package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	defaultPort = 8000
	defaultBind = map[int]Handler{
		400: DefaultHandle400(),
		401: DefaultHandle401(),
		403: DefaultHandle403(),
		404: DefaultHandle404(),
		405: DefaultHandle405(),
		500: DefaultHandle500(),
	}
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type Server struct {
	Server             *http.Server
	Port               int
	TLS                bool
	certFile, keyFile  string
	MaxRequestBodySize int
	Middle             *Middle
	Bind               *Bind
	group              *Group
	showLogo           bool
	showRouter         bool
}

// NewServer 新建http服务
func NewServer() *Server {
	server := &Server{
		Port:       defaultPort,
		TLS:        false,
		Middle:     &Middle{},
		Bind:       &Bind{m: defaultBind},
		group:      &Group{},
		showLogo:   true,
		showRouter: true,
	}
	server.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", defaultPort),
		Handler: http.HandlerFunc(server.ServeHTTP),
	}
	server.group.server = server
	return server
}

// Use 未实现
func (s *Server) Use(handler Handler) {
	s.Middle.Use(handler)
}

func (s *Server) ALL(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.ALL(pattern, handler)
	})
}

func (s *Server) GET(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.GET(pattern, handler)
	})
}

func (s *Server) POST(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.POST(pattern, handler)
	})
}

func (s *Server) PUT(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.PUT(pattern, handler)
	})
}

func (s *Server) DELETE(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.DELETE(pattern, handler)
	})
}

func (s *Server) PATCH(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.PATCH(pattern, handler)
	})
}

func (s *Server) HEAD(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.HEAD(pattern, handler)
	})
}

func (s *Server) CONNECT(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.CONNECT(pattern, handler)
	})
}

func (s *Server) OPTIONS(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.OPTIONS(pattern, handler)
	})
}

func (s *Server) TRACE(pattern string, handler Handler) *Server {
	return s.Group("", func(group *Group) {
		group.TRACE(pattern, handler)
	})
}

// Group 路由分组
func (s *Server) Group(pattern string, groups ...func(group *Group)) *Server {
	s.group.Group(pattern, groups...)
	return s
}

// SetPort 设置端口
func (s *Server) SetPort(port int) *Server {
	s.Port = port
	s.Server.Addr = fmt.Sprintf(":%d", port)
	return s
}

// SetMaxHeaderBytes
// 设置最大请求头大小,
// 例如:10<<20(10m)
// MaxHeaderBytes controls the maximum number of bytes the
// server will read parsing the request header's keys and
// values, including the request line. It does not limit the
// size of the request body.
// If zero, DefaultMaxHeaderBytes(1m) is used.
func (s *Server) SetMaxHeaderBytes(max int) {
	s.Server.MaxHeaderBytes = max
}

// SetTLS 设置tls
func (s *Server) SetTLS(certFile, keyFile string, tls ...bool) *Server {
	s.certFile = certFile
	s.keyFile = keyFile
	s.TLS = !(len(tls) > 0 && !tls[0])
	return s
}

// BindCodeHandler 绑定code对应handler
func (s *Server) BindCodeHandler(code int, handler Handler) *Server {
	s.Bind.Set(code, handler)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := newRequest(w, r)
	defer func() {
		if e := recover(); e != nil {
			msg := fmt.Sprint(e)
			if msg != MarkExit {
				request.ClearBody()
				request.WriteString(msg)
			}
		}
		request.done()
	}()
	defer func() {
		if e := recover(); e != nil {
			msg := fmt.Sprint(e)
			if msg != MarkExit {
				request.WriteErr(errors.New(msg))
			}
		}
		s.Bind.GetAndDo(request.GetStatusCode(), request)
	}()
	group := s.group.getGroup(request.URL.Path)
	if group == nil {
		request.SetStatusCode(404)
		return
	}
	group.Do(r.Method, request)
}

// Run 运行
func (s *Server) Run() (err error) {
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	if s.showLogo {
		fmt.Println(Logo)
	}
	if s.showRouter {
		s.printRouter()
	}
	log.Println("开启HTTP服务,端口:", s.Port)
	if !s.TLS {
		return s.Server.ListenAndServe()
	}
	return s.Server.ListenAndServeTLS(s.certFile, s.keyFile)
}

// ShowLogo 显示logo(INJOY)
func (s *Server) ShowLogo(b ...bool) {
	s.showLogo = !(len(b) > 0 && !b[0])
}

// ShowRouter 显示注册的路由
func (s *Server) ShowRouter(b ...bool) {
	s.showRouter = !(len(b) > 0 && !b[0])
}

func (s *Server) printRouter() {
	dealLength := func(s string, length int, fill string) string {
		for len(s) < length {
			s += fill
		}
		return s
	}
	maxApiUrl := 20
	maxApiPath := 20
	groupAll := s.group.getGroupAll()
	for _, v := range groupAll {
		if v.isBottom() && !v.route.IsEmpty() {
			if len(v.fullPath()) > maxApiUrl {
				maxApiUrl = len(v.fullPath())
			}
			for _, handler := range v.route.Clone() {
				handlerLen := len(handler.GetPath())
				if handlerLen > maxApiPath {
					maxApiPath = handlerLen
				}
			}
		}
	}
	maxLength := maxApiUrl + maxApiPath + 37
	print := func(method, fullPath, apiPath string) {
		fmt.Printf("|  %s|  %s|  %s|  %s|\n",
			dealLength(method, 10, " "),
			dealLength(s.Server.Addr, 10, " "),
			dealLength(fullPath, maxApiUrl+2, " "),
			dealLength(apiPath, maxApiPath+2, " "))
	}
	fmt.Println(dealLength("", maxLength, "="))
	fmt.Printf("|  %s|  %s|  %s|  %s|\n",
		dealLength("请求方式", 15, " "),
		dealLength("端口", 12, " "),
		dealLength("请求地址", maxApiUrl+7, " "),
		dealLength("接口路径", maxApiPath+7, " "))
	for _, v := range groupAll {
		if v.isBottom() && !v.route.IsEmpty() {
			fmt.Println(dealLength("", maxLength, "-"))
			for method, handler := range v.route.Clone() {
				print(method, v.fullPath(), handler.GetPath())
			}

		}
	}
	fmt.Println(dealLength("", maxLength, "="))
}
