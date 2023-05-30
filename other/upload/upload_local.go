package upload

import (
	"github.com/injoyai/base/bytes/crypt/md5"
	"io"
	"os"
	"path/filepath"
)

func NewLocal(dir string, rename ...bool) (Interface, error) {
	err := os.MkdirAll(dir, 0777)
	return &Local{
		dir:    dir,
		rename: len(rename) > 0 && rename[0],
	}, err
}

type Local struct {
	dir    string
	rename bool
}

func (this *Local) Save(filename string, reader io.Reader) (string, error) {
	if this.rename {
		filename = md5.Encrypt(filename)
	}
	filename = filepath.Join(this.dir, filename)
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return filename, err
}
