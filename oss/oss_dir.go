package oss

import (
	"io"
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

// RangeFileInfo 遍历目录
func RangeFileInfo(dir string, fn func(info *FileInfo) (bool, error), level ...int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		next, err := fn(&FileInfo{
			FileInfo: info,
			Dir:      dir,
		})
		if err != nil {
			return err
		}
		if !next {
			break
		}
		if len(level) > 0 && level[0] != 0 && entry.IsDir() {
			//大于0(限制多少层)或者小于0(无限制层数),表示继续往下遍历
			if err := RangeFileInfo(filepath.Join(dir, info.Name()), fn, level[0]-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadFileInfos 获取目录下的所有文件信息(包括文件夹)
func ReadFileInfos(dir string, level ...int) ([]os.FileInfo, error) {
	files := []os.FileInfo(nil)
	err := RangeFileInfo(dir, func(info *FileInfo) (bool, error) {
		files = append(files, info)
		return true, nil
	}, level...)
	return files, err
}

// ReadFilenames 获取目录下的所有文件名称
func ReadFilenames(dir string, level ...int) ([]string, error) {
	filenames := []string(nil)
	err := RangeFileInfo(dir, func(info *FileInfo) (bool, error) {
		if !info.IsDir() {
			filenames = append(filenames, info.FullName())
		}
		return true, nil
	}, level...)
	return filenames, err
}

// ReadDirNames 获取目录下的所有目录
func ReadDirNames(dir string, level ...int) ([]string, error) {
	dirNames := []string(nil)
	err := RangeFileInfo(dir, func(info *FileInfo) (bool, error) {
		if info.IsDir() {
			dirNames = append(dirNames, info.FullName())
		}
		return true, nil
	}, level...)
	return dirNames, err
}

// RangeFile 遍历目录的文件
func RangeFile(dir string, fn func(info *FileInfo, f *os.File) (bool, error), level ...int) error {
	return RangeFileInfo(dir, func(info *FileInfo) (bool, error) {
		if !info.IsDir() {
			f, err := os.Open(filepath.Join(dir, info.Name()))
			if err != nil {
				return false, err
			}
			defer f.Close()
			return fn(info, f)
		}
		return true, nil
	}, level...)
}

// RangeFileBytes 遍历目录的文件字节
func RangeFileBytes(dir string, fn func(info *FileInfo, bs []byte) bool, level ...int) error {
	return RangeFile(dir, func(info *FileInfo, f *os.File) (bool, error) {
		bs, err := io.ReadAll(f)
		if err != nil {
			return false, err
		}
		return fn(info, bs), nil
	}, level...)
}

type FileInfo struct {
	os.FileInfo
	Dir string
}

func (this *FileInfo) FullName() string {
	return filepath.Join(this.Dir, this.FileInfo.Name())
}

func (this *FileInfo) Filename() string {
	return filepath.Join(this.Dir, this.FileInfo.Name())
}
