package upload

import (
	"github.com/injoyai/base/bytes/crypt/md5"
	"github.com/injoyai/goutil/oss"
	"io"
	"os"
	"path/filepath"
)

var DefaultLocal = NewLocal("./data/upload/")

func NewLocal(dir string, rename ...bool) Interface {
	return &Local{
		dir:    dir,
		rename: len(rename) > 0 && rename[0],
	}
}

type Local struct {
	dir    string //保存的目录
	rename bool   //是否重命名
}

func (this *Local) Save(name string, reader io.Reader) (string, error) {
	dir := this.dir
	dir, name = filepath.Split(name)
	dir = filepath.Join(this.dir, dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	if this.rename {
		name = md5.Encrypt(name)
	}
	filename := filepath.Join(dir, name)
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if reader != nil {
		_, err = io.Copy(f, reader)
	}
	return filename, err
}

func (this *Local) List() ([]string, error) {
	return oss.ReadFilenames(this.dir, -1)
}
