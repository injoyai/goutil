package oss

import (
	"os"
	"runtime"
)

// Exists 是否存在
func Exists(name string) bool {
	stat, err := os.Stat(name)
	return stat != nil && !os.IsNotExist(err)
}

// Stat 获取文件信息
func Stat(filename string) (os.FileInfo, bool, error) {
	stat, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, false, err
	} else if err != nil && os.IsNotExist(err) {
		return nil, false, nil
	}
	return stat, true, nil
}

// IsDir 是否是文件夹
func IsDir(name string) bool {
	s, err := os.Stat(name)
	return err == nil && s.IsDir()
}

// IsFile 是否是文件
func IsFile(name string) bool {
	s, err := os.Stat(name)
	return err == nil && !s.IsDir()
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}
