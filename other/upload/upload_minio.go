package upload

import (
	"fmt"
	"github.com/injoyai/base/crypt/md5"
	"github.com/minio/minio-go"
	"io"
	"strings"
)

func NewMinio(cfg *MinioConfig) (Interface, error) {
	secure := false
	endpoint := cfg.Endpoint
	switch {
	case strings.HasPrefix(cfg.Endpoint, "https://"):
		secure = true
		endpoint = strings.TrimPrefix(cfg.Endpoint, "https://")
	case strings.HasPrefix(cfg.Endpoint, "http://"):
		endpoint = strings.TrimPrefix(cfg.Endpoint, "http://")
	default:
		endpoint = cfg.Endpoint
		cfg.Endpoint = "http://" + cfg.Endpoint
	}
	cli, err := minio.New(endpoint, cfg.AccessKey, cfg.SecretKey, secure)
	if err != nil {
		return nil, err
	}
	return &Minio{Client: cli, cfg: cfg}, nil
}

type MinioConfig struct {
	Endpoint   string //地址 http://127.0.0.1:9000
	AccessKey  string //访问key
	SecretKey  string //秘钥
	BucketName string //桶名称
	Rename     bool   //是否重命名文件
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

func (this *Minio) List() ([]string, error) {
	list := []string(nil)
	for v := range this.Client.ListObjects(this.cfg.BucketName, "", false, make(chan struct{})) {
		list = append(list, fmt.Sprintf("%s/%s/%s", this.cfg.Endpoint, this.cfg.BucketName, v.Key))
	}
	return list, nil
}
