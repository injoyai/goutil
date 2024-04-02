package upload

import "io"

const (
	TypeMinio = "minio"
	TypeLocal = "local"
	TypeQiniu = "qiniu"
)

type Interface interface {
	Save(filename string, reader io.Reader) (string, error)
	List() ([]string, error)
}
