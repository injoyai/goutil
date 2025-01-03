package upload

import (
	"io"
)

var _ Uploader = (*HTTP)(nil)

type HTTP struct {
}

func (H HTTP) Upload(filename string, reader io.Reader) (URL, error) {
	//TODO implement me
	panic("implement me")
}

func (H HTTP) List(join ...string) ([]*Info, error) {
	//TODO implement me
	panic("implement me")
}
