package disk

import (
	"errors"
	"github.com/injoyai/goutil/oss"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/encrypt"
	"io/fs"
	"strings"
)

var _ Disker = (*Minio)(nil)

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

func (this *Minio) List(dir string) ([]fs.FileInfo, error) {
	list := []fs.FileInfo(nil)
	for v := range this.Client.ListObjectsV2(this.cfg.BucketName, dir, true, make(chan struct{})) {
		list = append(list, &FileInfo{
			name:  v.Key,
			size:  v.Size,
			isDir: v.Size == 0 && v.Owner.ID == "",
			time:  v.LastModified,
		})
	}
	return list, nil
}

func (this *Minio) Rename(filename, newName string) error {
	dest, err := minio.NewDestinationInfo(this.cfg.BucketName, filename, encrypt.NewSSE(), nil)
	if err != nil {
		return err
	}
	sour := minio.NewSourceInfo(this.cfg.BucketName, newName, encrypt.NewSSE())
	if err := this.Client.CopyObject(dest, sour); err != nil {
		return err
	}
	return this.Client.RemoveObject(this.cfg.BucketName, filename)
}

func (this *Minio) Mkdir(dir string) error {
	return errors.New("未实现")
}

func (this *Minio) Upload(filename string, f fs.File) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	_, err = this.Client.PutObject(this.cfg.BucketName, filename, f, fi.Size(), minio.PutObjectOptions{})
	return err
}

func (this *Minio) Download(filename, localFilename string) error {
	ob, err := this.Client.GetObject(this.cfg.BucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer ob.Close()
	return oss.New(localFilename, ob)
}

func (this *Minio) Delete(filename string) error {
	return this.Client.RemoveObject(this.cfg.BucketName, filename)
}
