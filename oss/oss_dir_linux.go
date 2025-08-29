package oss

import (
	"path/filepath"
)

// UserDataDir 系统用户数据路径,AppData/Local
func UserDataDir(join ...string) string {
	return filepath.Join(append([]string{"/var/lib/"}, join...)...)
}

// UserStartupDir 自启路径,这个对linux没啥用,需要编辑文件,而不是放到目录就行
func UserStartupDir(join ...string) string {
	return filepath.Join(append([]string{"~/.config/autostart/"}, join...)...)
}
