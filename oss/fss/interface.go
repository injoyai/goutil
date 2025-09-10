package fss

import (
	"io"
	"io/fs"
)

type Root interface {
	io.Closer                                     //
	Open(filename string) (fs.File, error)        //打开文件
	Create(filename string) (fs.File, error)      //新建文件
	Remove(filename string) error                 //删除文件/目录
	Stat(filename string) (fs.FileInfo, error)    //获取文件信息
	Rename(oldFilename, newFilename string) error //重命名/移动
	ReadDir(dir string) ([]fs.FileInfo, error)    //读取目录
	Mkdir(name string, perm fs.FileMode) error    //创建目录
}
