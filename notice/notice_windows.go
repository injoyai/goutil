//go:build windows
// +build windows

package notice

import (
	"github.com/go-toast/toast"
	"syscall"
	"unsafe"
)

func NewWindows() Interface { return &Windows{} }

// Windows 右下角通知和弹窗
type Windows struct{}

func (this *Windows) Publish(msg *Message) error {
	switch msg.Target {
	case TargetPopup, TargetPop:

		//弹窗通知,会阻塞,等待用户关掉才能返回
		user32dll, err := syscall.LoadLibrary("user32.dll")
		if err != nil {
			return err
		}
		user32 := syscall.NewLazyDLL("user32.dll")
		MessageBoxW := user32.NewProc("MessageBoxW")
		_, _, err = MessageBoxW.Call(
			uintptr(0),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msg.Content))),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msg.Title))),
			uintptr(0),
		)
		defer syscall.FreeLibrary(user32dll)
		if err != nil && err.Error() == "The operation completed successfully." {
			return nil
		}
		return err

	default:

		//右下角通知
		appID := msg.Tag.Conv().GetString("AppID", "Microsoft.Windows.Shell.RunDialog")
		notification := toast.Notification{
			AppID:    appID,
			Title:    msg.Title,
			Message:  msg.Content,
			Audio:    toast.Default,
			Duration: toast.Short,
		}
		return notification.Push()

	}
}