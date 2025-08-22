package upload

import (
	"fmt"
	"github.com/injoyai/conv"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Uploader interface {
	Upload(filename string, reader io.Reader) (URL, error) //上传文件,返回下载地址
	List(join ...string) ([]*Info, error)                  //目录列表
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
	TypeLocal    = "local"
	TypeMinio    = "minio"
	TypeBaidu    = "baidu"
	TypeCloud189 = "cloud189"
	TypeFtp      = "ftp"
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
	case TypeFtp:
		return DialFTP(
			cfg.GetString("address"),
			cfg.GetString("username"),
			cfg.GetString("password"),
		)
	case TypeBaidu:
		return NewBaidu(), nil
	case TypeCloud189:
		return NewCloud189(Cloud189Config{
			Username: cfg.GetString("username"),
			Password: cfg.GetString("password"),
			Size:     cfg.GetInt("size"),
			Retry:    cfg.GetInt("retry"),
		})
	}
	return nil, fmt.Errorf("未知类型:%s", Type)
}

// SyncDir 同步目录
func SyncDir(up Uploader, localDir, remoteDir string) error {
	entries, err := os.ReadDir(localDir)
	if err != nil {
		return err
	}
	for _, info := range entries {
		if info.IsDir() {
			newRemoteDir := filepath.Join(remoteDir, info.Name())
			newRemoteDir = strings.ReplaceAll(newRemoteDir, "\\", "/")
			if err := SyncDir(up, filepath.Join(localDir, info.Name()), newRemoteDir); err != nil {
				return err
			}
		} else {
			f, err := os.Open(filepath.Join(localDir, info.Name()))
			if err != nil {
				return err
			}
			remoteFilename := filepath.Join(remoteDir, info.Name())
			remoteFilename = strings.ReplaceAll(remoteFilename, "\\", "/")
			if _, err := up.Upload(remoteFilename, f); err != nil {
				return err
			}
		}
	}
	return nil
}
