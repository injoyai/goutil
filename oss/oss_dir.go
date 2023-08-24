package oss

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultName = "injoy"
)

// ExecName 当前执行的程序名称
func ExecName() string {
	fullName, _ := os.Executable()
	return fullName
}

// ExecDir 当前执行的程序路径
func ExecDir() string {
	return filepath.Dir(ExecName())
}

// FuncName 当前执行的函数名称
func FuncName() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

// FuncDir 当前执行的函数路径
func FuncDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

// UserDir 系统用户路径
func UserDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}

// UserDataDir 系统用户数据路径
func UserDataDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Local"}, join...)...)
}

// UserLocalDir 系统用户本地数据路径
func UserLocalDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Local"}, join...)...)
}

// UserStartupDir 自启路径
func UserStartupDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup"}, join...)...)
}

func UserInjoyDir(join ...string) string {
	return UserLocalDir(append([]string{DefaultName}, join...)...)
}

// UserDefaultDir 默认系统用户数据子路径(个人使用)
func UserDefaultDir(join ...string) string {
	return UserInjoyDir(join...)
}

/*



 */

// NewDir 新建文件夹
// @path,路径
func NewDir(path string) error {
	return os.MkdirAll(path, defaultPerm)
}

func DelDir(dir string) error {
	return os.RemoveAll(dir)
}

// ReadDirFunc 遍历目录
func ReadDirFunc(dir string, fn func(info os.FileInfo) error) error {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, info := range fileInfos {
		if err = fn(info); err != nil {
			return err
		}
	}
	return nil
}

// ReadFileInfos 获取目录下的所有文件信息
func ReadFileInfos(dir string) ([]os.FileInfo, error) {
	files := []os.FileInfo(nil)
	err := ReadDirFunc(dir, func(info os.FileInfo) error {
		if !info.IsDir() {
			files = append(files, info)
		}
		return nil
	})
	return files, err
}

// ReadFilenames 获取目录下的所有文件名称
func ReadFilenames(dir string) ([]string, error) {
	filenames := []string(nil)
	err := ReadDirFunc(dir, func(info os.FileInfo) error {
		filenames = append(filenames, filepath.Join(dir, info.Name()))
		return nil
	})
	return filenames, err
}

// ReadDirNames 获取目录下的所有目录
func ReadDirNames(dir string) ([]string, error) {
	dirNames := []string(nil)
	err := ReadDirFunc(dir, func(info os.FileInfo) error {
		if info.IsDir() {
			dirNames = append(dirNames, filepath.Join(dir, info.Name()))
		}
		return nil
	})
	return dirNames, err
}
