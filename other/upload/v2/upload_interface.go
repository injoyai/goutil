package upload

import (
	"fmt"
	"github.com/injoyai/conv"
	"io"
)

type Uploader interface {
	Upload(filename string, reader io.Reader) (URL, error) //上传文件,返回下载地址
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
	Rename Equal = 2 //重命名文件
)

type Equal uint8

var (
	TypeLocal = "local"
	TypeMinio = "minio"
	TypeQiniu = "qiniu"
	TypeBaidu = "baidu"
)

func New(Type string, cfg conv.Extend) (Uploader, error) {
	switch Type {
	case TypeLocal:
		return NewLocal(cfg.GetString("dir")), nil
	case TypeMinio:
		return NewMinio(&MinioConfig{
			Endpoint:   cfg.GetString("endpoint"),
			AccessKey:  cfg.GetString("accessKey"),
			SecretKey:  cfg.GetString("secretKey"),
			BucketName: cfg.GetString("bucketName"),
		})
	case TypeQiniu:
		return NewQiniu(&QiniuConfig{
			Key:    cfg.GetString("key"),
			Secret: cfg.GetString("secret"),
			Domain: cfg.GetString("domain"),
			Space:  cfg.GetString("space"),
		})
	case TypeBaidu:
		return NewBaidu(), nil
	}
	return nil, fmt.Errorf("未知类型:%s", Type)
}
