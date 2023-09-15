package http_handler

import (
	"github.com/injoyai/goutil/oss"
	"net/http"
	"path/filepath"
	"strings"
)

type Option func(w http.ResponseWriter, r *http.Request) bool

var DefaultFile = NewFile("", "")

type File struct {
	Prefix  string   //前缀
	BootDir string   //根目录
	Option  []Option //中间件
}

func (this *File) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, v := range this.Option {
		if !v(w, r) {
			return
		}
	}
	filename := r.URL.Path
	if !strings.HasPrefix(r.URL.Path, this.Prefix) {
		http.NotFound(w, r)
		return
	}
	filename = filepath.Join(this.BootDir, filename[len(this.Prefix):])
	if !oss.Exists(filename) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, filename)
}

func NewFile(dir string, prefix string, options ...Option) http.Handler {
	if len(dir) == 0 {
		dir = "."
	}
	return &File{
		Prefix:  prefix,
		BootDir: dir,
		Option:  options,
	}
}
