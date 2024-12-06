package upload

import (
	"fmt"
	"github.com/minio/minio-go"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func NewMinio(cfg *MinioConfig) (*Minio, error) {
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
}

type Minio struct {
	cfg *MinioConfig
	*minio.Client
}

func (this *Minio) Upload(filename string, reader io.Reader) (URL, error) {
	_, err := this.PutObject(this.cfg.BucketName, filename, reader, -1, minio.PutObjectOptions{})
	return HttpUrl(fmt.Sprintf("%s/%s/%s", this.cfg.Endpoint, this.cfg.BucketName, filename)), err
}

func (this *Minio) List(join ...string) ([]*Info, error) {
	list := []*Info(nil)
	for v := range this.Client.ListObjectsV2(this.cfg.BucketName, filepath.Join(join...), true, make(chan struct{})) {
		list = append(list, &Info{
			Name: fmt.Sprintf("%s/%s/%s", this.cfg.Endpoint, this.cfg.BucketName, v.Key),
			Size: v.Size,
			Dir:  v.Size == 0 && v.Owner.ID == "",
			Time: v.LastModified.Unix(),
		})
	}
	return list, nil
}

type HttpUrl string

func (this HttpUrl) String() string {
	return string(this)
}

func (this HttpUrl) Download(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	req, err := http.NewRequest("GET", this.String(), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
