package pprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

/*
Run

网页打开
http://localhost:6060/debug/pprof/

内存使用情况
go tool pprof http://localhost:6060/debug/pprof/heap

cpu使用情况
go tool pprof http://localhost:6060/debug/pprof/profile?scends=30
*/
func Run(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
