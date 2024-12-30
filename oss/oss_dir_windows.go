package oss

import (
	"os"
	"path/filepath"
)

// UserDataDir 系统用户数据路径,AppData/Local
func UserDataDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Local"}, join...)...)
}

// UserStartupDir 自启路径
func UserStartupDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup"}, join...)...)
}
