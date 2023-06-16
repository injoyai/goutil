package oss

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultName = "injoy"
)

func ExecName() string {
	fullName, _ := os.Executable()
	return fullName
}

func ExecDir() string {
	return filepath.Dir(ExecName())
}

func FuncName() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

func FuncDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func UserDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}

func UserDataDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData"}, join...)...)
}

func UserLocalDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Local"}, join...)...)
}

func UserDefaultDir() string {
	return UserLocalDir(DefaultName)
}

func Exists(name string) bool {
	stat, err := os.Stat(name)
	return stat != nil && !os.IsNotExist(err)
}
