package http_handler

import (
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/str"
	"github.com/injoyai/logs"
	"net/http"
	"path/filepath"
)

var DefaultFile http.Handler

func init() {
	DefaultFile = NewFile("", "")
}

type File struct {
	Prefix  string         //前缀
	BootDir string         //根目录
	Option  []http.Handler //中间件
}

func (this *File) pageNotFind(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte("404 page not find"))
}

func (this *File) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := str.CropFirst(r.URL.Path, this.Prefix, false)
	filename = filepath.Join(this.BootDir, filename)
	logs.Debug(filename)
	if !oss.Exists(filename) {
		w.WriteHeader(404)
		w.Write([]byte("404 page not find"))
		return
	}
	bs, err := oss.ReadBytes(filename)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404 page not find"))
		return
	}
	w.WriteHeader(200)
	w.Write(bs)
}

func NewFile(dir string, prefix string, options ...http.Handler) http.Handler {
	if len(dir) == 0 {
		dir = "."
	}
	return &File{
		Prefix:  prefix,
		BootDir: dir,
		Option:  options,
	}
}
