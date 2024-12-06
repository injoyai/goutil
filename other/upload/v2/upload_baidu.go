package upload

import (
	"errors"
	"io"
)

var _ Uploader = (*Baidu)(nil)

func NewBaidu() *Baidu {
	return &Baidu{}
}

type Baidu struct {
}

func (this *Baidu) Upload(filename string, reader io.Reader) (URL, error) {
	return nil, errors.New("未实现")
}

func (this *Baidu) List(join ...string) ([]*Info, error) {
	return nil, errors.New("未实现")
}
