package upload

import (
	"fmt"
	"github.com/injoyai/base/bytes/crypt/md5"
	"github.com/minio/minio-go"
	"io"
)

func NewMinio(cfg *MinioConfig) (Interface, error) {
	cli, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, false)
	if err != nil {
		return nil, err
	}
	return &Minio{Client: cli, cfg: cfg}, nil
}

type MinioConfig struct {
	Endpoint   string //地址
	AccessKey  string //访问key
	SecretKey  string //秘钥
	BucketName string //桶名称
	Rename     bool   //重命名
}

type Minio struct {
	cfg *MinioConfig
	*minio.Client
}

func (this *Minio) Save(filename string, reader io.Reader) (string, error) {
	if this.cfg.Rename {
		filename = md5.Encrypt(filename)
	}
	_, err := this.PutObject(this.cfg.BucketName, filename, reader, -1, minio.PutObjectOptions{})
	return fmt.Sprintf("%s/%s/%s", this.cfg.Endpoint, this.cfg.BucketName, filename), err
}
