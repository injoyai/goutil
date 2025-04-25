package disk

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Disker interface {
	List(dir string) ([]fs.FileInfo, error)    //目录
	Rename(filename, newName string) error     //重命名
	Mkdir(dir string) error                    //创建目录
	Upload(filename string, f fs.File) error   //上传
	Download(filename, localName string) error //下载
	Delete(filename string) error              //删除
}

// UploadDir 上传目录
func UploadDir(disk Disker, localDir, remoteDir string) error {
	entries, err := os.ReadDir(localDir)
	if err != nil {
		return err
	}
	for _, info := range entries {
		if info.IsDir() {
			newRemoteDir := filepath.Join(remoteDir, info.Name())
			newRemoteDir = strings.ReplaceAll(newRemoteDir, "\\", "/")
			if err := UploadDir(disk, filepath.Join(localDir, info.Name()), newRemoteDir); err != nil {
				return err
			}
		} else {
			f, err := os.Open(filepath.Join(localDir, info.Name()))
			if err != nil {
				return err
			}
			remoteFilename := filepath.Join(remoteDir, info.Name())
			remoteFilename = strings.ReplaceAll(remoteFilename, "\\", "/")
			if err = disk.Upload(remoteFilename, f); err != nil {
				return err
			}
		}
	}
	return nil
}

func SyncDir(disk Disker, localDir, remoteDir string) error {
	return nil
}

var _ fs.FileInfo = (*FileInfo)(nil)

type FileInfo struct {
	name  string
	size  int64
	mode  fs.FileMode
	time  time.Time
	isDir bool
}

func (this *FileInfo) Name() string {
	return this.name
}

func (this *FileInfo) Size() int64 {
	return this.size
}

func (this *FileInfo) Mode() fs.FileMode {
	return this.mode
}

func (this *FileInfo) ModTime() time.Time {
	return this.time
}

func (this *FileInfo) IsDir() bool {
	return this.isDir
}

func (this *FileInfo) Sys() any {
	return nil
}
