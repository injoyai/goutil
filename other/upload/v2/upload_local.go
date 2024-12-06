package upload

import (
	"github.com/injoyai/goutil/oss"
	"io"
	"os"
	"path/filepath"
)

var _ Uploader = (*Local)(nil)

var DefaultDir = oss.ExecDir("/data/upload/")

func NewLocal(dir string) *Local {
	return &Local{
		dir: dir,
	}
}

type Local struct {
	dir string
}

func (this *Local) Upload(filename string, reader io.Reader) (URL, error) {
	dir, name := filepath.Split(filename)
	dir = filepath.Join(this.dir, dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	filename = filepath.Join(dir, name)
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if reader != nil {
		_, err = io.Copy(f, reader)
	}
	return LocalPath(filename), err
}

func (this *Local) List(join ...string) ([]*Info, error) {
	dir := filepath.Join(this.dir, filepath.Join(join...))
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	ls := make([]*Info, len(dirs))
	for i := range dirs {
		info, err := dirs[i].Info()
		if err != nil {
			return nil, err
		}
		ls[i] = &Info{
			Name: info.Name(),
			Size: info.Size(),
			Dir:  info.IsDir(),
			Time: info.ModTime().Unix(),
		}
	}
	return ls, nil
}

type LocalPath string

func (this LocalPath) String() string {
	return string(this)
}

func (this LocalPath) Download(filename string) error {
	f, err := os.Open(this.String())
	if err != nil {
		return err
	}
	defer f.Close()
	f2, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f2.Close()
	io.Copy(f2, f)
	return nil
}
