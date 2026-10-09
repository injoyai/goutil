package oss

import (
	"os"
	"path/filepath"
)

// UserDataDir 系统用户数据路径, ~/Library/Application Support
func UserDataDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "Library/Application Support"}, join...)...)
}

// UserStartupDir 自启路径, ~/Library/LaunchAgents
func UserStartupDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "Library/LaunchAgents"}, join...)...)
}
