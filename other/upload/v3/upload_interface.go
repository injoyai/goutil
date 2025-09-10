package upload

import (
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss/fss"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Uploader interface {
	fss.Root
	Upload(filename string, reader io.Reader) (string, error) //上传文件,返回下载地址
	Download(filename, localFilename string) error            //下载文件
	Dir(join ...string) ([]*Info, error)                      //目录列表
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
	TypeSmb      = "smb"
)

func New(Type string, cfg conv.Extend) (Uploader, error) {
	switch Type {
	case TypeSmb:
		return NewSmb(&SmbConfig{
			Host:      cfg.GetString("host"),
			Username:  cfg.GetString("username"),
			Password:  cfg.GetString("password"),
			ShareName: cfg.GetString("shareName"),
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
