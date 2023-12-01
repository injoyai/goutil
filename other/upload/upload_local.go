package upload

import (
	"github.com/injoyai/base/bytes/crypt/md5"
	"io"
	"os"
	"path/filepath"
)

var DefaultLocal = NewLocal()

func NewLocal(rename ...bool) Interface {
	return &Local{rename: len(rename) > 0 && rename[0]}
}

type Local struct {
	rename bool
}

func (this *Local) Save(filename string, reader io.Reader) (string, error) {
	dir, name := filepath.Split(filename)
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return "", err
	}
	if this.rename {
		filename = filepath.Join(dir, md5.Encrypt(name))
	}
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
