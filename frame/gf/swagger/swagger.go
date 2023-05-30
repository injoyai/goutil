package swagger

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/injoyai/conv"
	"io/ioutil"
	"net/http"
)

func Handler(jsonPath ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := conv.GetDefaultString("/swagger/swagger.json", jsonPath...)
		w.Write([]byte(fmt.Sprintf(ui, path)))
		w.WriteHeader(200)
	})
}

type Swagger struct {
	Prefix string //路由前缀
	Path   string //文件路径
}

func (this *Swagger) Name() string {
	return "swagger"
}

func (this *Swagger) Author() string {
	return "injoy"
}

func (this *Swagger) Version() string {
	return "v2.0"
}

func (this *Swagger) Description() string {
	return ""
}

func (this *Swagger) Install(s *ghttp.Server) error {
	s.Group(this.Prefix, func(group *ghttp.RouterGroup) {
		group.GET("/", func(r *ghttp.Request) {
			r.Response.Write([]byte(fmt.Sprintf(ui, this.Prefix+"/swagger.json")))
			r.ExitAll()
		})
		group.GET("/swagger.json", func(r *ghttp.Request) {
			bytes, err := ioutil.ReadFile(this.Path)
			if err == nil {
				r.Response.ClearBuffer()
				r.Response.WriteExit(bytes)
				r.ExitAll()
			}
		})
	})
	return nil
}

func (this *Swagger) Remove() error {
	return nil
}
