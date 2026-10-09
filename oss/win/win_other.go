//go:build !windows
// +build !windows

package win

import (
	"errors"
)

// CreateStartupShortcut 创建自启快捷方式, 仅windows支持
func CreateStartupShortcut(target string) error {
	return errors.New("创建自启快捷方式仅支持windows")
}
