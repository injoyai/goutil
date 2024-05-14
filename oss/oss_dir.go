package oss

import (
	"github.com/injoyai/conv"
	"io"
	"io/fs"
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

// UserDir 系统用户路径 C:\Users\QL1211
func UserDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}

// UserDataDir 系统用户数据路径,AppData/Local
func UserDataDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Local"}, join...)...)
}

// UserLocalDir 系统用户本地数据路径,AppData/Local
func UserLocalDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Local"}, join...)...)
}

// UserStartupDir 自启路径
func UserStartupDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup"}, join...)...)
}

func UserDesktopDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "Desktop"}, join...)...)
}

// UserInjoyDir 个人数据路径
func UserInjoyDir(join ...string) string {
	return UserLocalDir(append([]string{DefaultName}, join...)...)
}

// UserInjoyCacheDir 个人缓存数据路径
func UserInjoyCacheDir(join ...string) string {
	return UserInjoyDir(append([]string{"/data/cache"}, join...)...)
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

// ReadFileInfos 获取目录下的所有文件信息(包括文件夹)
func ReadFileInfos(dir string) ([]os.FileInfo, error) {
	files := []os.FileInfo(nil)
	err := RangeFileInfo(dir, func(info os.FileInfo) (bool, error) {
		files = append(files, info)
		return true, nil
	})
	return files, err
}

// ReadFilenames 获取目录下的所有文件名称
func ReadFilenames(dir string, levels ...int) ([]string, error) {
	level := conv.DefaultInt(0, levels...)
	filenames := []string(nil)
	err := RangeFileInfo(dir, func(info os.FileInfo) (bool, error) {
		if info.IsDir() {
			if level != 0 {
				if level > 0 {
					level--
				}
				child, err := ReadFilenames(filepath.Join(dir, info.Name()), level)
				if err != nil {
					return false, err
				}
				filenames = append(filenames, child...)
			}
		} else {
			filenames = append(filenames, filepath.Join(dir, info.Name()))
		}
		return true, nil
	})
	return filenames, err
}

// ReadDirNames 获取目录下的所有目录
func ReadDirNames(dir string) ([]string, error) {
	dirNames := []string(nil)
	err := RangeFileInfo(dir, func(info os.FileInfo) (bool, error) {
		if info.IsDir() {
			dirNames = append(dirNames, filepath.Join(dir, info.Name()))
		}
		return true, nil
	})
	return dirNames, err
}

// RangeFileInfo 遍历目录
func RangeFileInfo(dir string, fn func(info fs.FileInfo) (bool, error)) error {
	entrys, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entrys {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		next, err := fn(info)
		if err != nil {
			return err
		}
		if !next {
			break
		}
	}
	return nil
}

// RangeFile 遍历目录的文件
func RangeFile(dir string, fn func(info os.FileInfo, f *os.File) (bool, error)) error {
	return RangeFileInfo(dir, func(info os.FileInfo) (bool, error) {
		if !info.IsDir() {
			f, err := os.Open(filepath.Join(dir, info.Name()))
			if err != nil {
				return false, err
			}
			defer f.Close()
			return fn(info, f)
		}
		return true, nil
	})
}

// RangeFileBytes 遍历目录的文件字节
func RangeFileBytes(dir string, fn func(info os.FileInfo, bs []byte) bool) error {
	return RangeFile(dir, func(info os.FileInfo, f *os.File) (bool, error) {
		bs, err := io.ReadAll(f)
		if err != nil {
			return false, err
		}
		return fn(info, bs), nil
	})
}
