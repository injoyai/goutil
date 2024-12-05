package upload

import (
	"github.com/injoyai/conv"
	"io"
)

type Uploader interface {
	Upload(filename string, reader io.Reader) (URL, error) //保存文件,返回下载地址
	List(join ...string) ([]*Info, error)
}

type URL interface {
	String() string
	Download(filename string) error
}

type Info struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Dir  bool   `json:"dir"`
	Time int64  `json:"time"`
}

const (
	Cover  Equal = 0 //覆盖老文件
	Keep   Equal = 1 //保持旧文件
	Rename Equal = 2 //重命名为随机
)

type Equal uint8

var (
	TypeLocal = "local"
	TypeMinio = "minio"
)

func New(Type string, cfg conv.Extend) Uploader {
	switch Type {
	case TypeLocal:
		return NewLocal(
			cfg.GetString("dir"),
			cfg.GetBool("rename"),
		)
	}
	return nil
}
