//go:build windows
// +build windows

package win

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/injoyai/goutil/oss"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Shortcut the shortcut (.lnk file) property struct
type Shortcut struct {
	// Shortcut (.lnk file) path
	ShortcutPath string
	// Shortcut target: a file path or a website
	Target string
	// Shortcut icon path, default: "%SystemRoot%\\System32\\SHELL32.dll,0"
	IconLocation string
	// Arguments of shortcut
	Arguments string
	// Description of shortcut
	Description string
	// Hotkey of shortcut
	Hotkey string
	// WindowStyle, "1"(default) for default size and location; "3" for maximized window; "7" for minimized window
	WindowStyle string
	// Working directory of shortcut
	WorkingDirectory string
}

// Create 创建快捷方式
func (this *Shortcut) Create() error {
	return CreateShortcut(this)
}

// Remove 移除快捷方式
func (this *Shortcut) Remove() error {
	return os.Remove(this.ShortcutPath)
}

// NewStartupShortcut 新建开机自启实例
func NewStartupShortcut(target string) *Shortcut {
	_, name := filepath.Split(target)
	name = strings.Split(name, ".")[0]
	shortcutPath := oss.UserStartupDir(name + ".lnk")
	return &Shortcut{
		ShortcutPath: shortcutPath,
		Target:       target,
		WindowStyle:  "1",
	}
}

// NewDesktopShortcut 新建桌面快捷方式实例
func NewDesktopShortcut(name, target string) *Shortcut {
	u, err := user.Current()
	if err != nil {
		return nil
	}
	shortcutPath := filepath.Join(u.HomeDir, "Desktop", name+".lnk")
	return &Shortcut{
		ShortcutPath: shortcutPath,
		Target:       target,
		WindowStyle:  "1",
	}
}

// CreateStartupShortcut 创建开机自启快捷方式
func CreateStartupShortcut(target string) error {
	return NewStartupShortcut(target).Create()
}

// RemoveStartupShortcut 一处开机自启快捷方式
func RemoveStartupShortcut(target string) error {
	return NewStartupShortcut(target).Remove()
}

// CreateDesktopShortcut 创建桌面快捷方式
// 例 CreateDesktopShortcut("google","https://google.cn")
func CreateDesktopShortcut(name, target string) error {
	return NewDesktopShortcut(name, target).Create()
}

// RemoveDesktopShortcut 删除桌面快捷方式
func RemoveDesktopShortcut(name string) error {
	return NewDesktopShortcut(name, "").Remove()
}

// CreateShortcut 创建快捷方式
func CreateShortcut(shortcut *Shortcut) error {
	if shortcut.IconLocation == "" {
		shortcut.IconLocation = "%SystemRoot%\\System32\\SHELL32.dll,0"
	}
	if shortcut.WindowStyle == "" {
		shortcut.WindowStyle = "1"
	}
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcut.ShortcutPath)
	if err != nil {
		return err
	}

	idispatch := cs.ToIDispatch()
	_, err = oleutil.PutProperty(idispatch, "IconLocation", shortcut.IconLocation)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "TargetPath", shortcut.Target)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "Arguments", shortcut.Arguments)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "Description", shortcut.Description)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "Hotkey", shortcut.Hotkey)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "WindowStyle", shortcut.WindowStyle)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(idispatch, "WorkingDirectory", shortcut.WorkingDirectory)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(idispatch, "Save")
	if err != nil {
		return err
	}
	return nil
}
