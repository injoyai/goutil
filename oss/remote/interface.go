package remote

import (
	"io"
	"io/fs"
)

type OS interface {
	Dir() ([]fs.FileInfo, error)                    //查看目录
	Open(name string) (io.ReadWriteCloser, error)   //打开文件
	Remove(name string) error                       //移除文件
	Create(name string) (io.ReadWriteCloser, error) //创建文件
}

type Info struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Dir  bool   `json:"dir"`
	Time int64  `json:"time"`
}
