package shell

import (
	"errors"
	"fmt"
	"github.com/injoyai/goutil/oss/linux/bash"
	"github.com/injoyai/goutil/oss/linux/systemctl"
	"github.com/injoyai/goutil/str"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Execf(format string, args ...interface{}) (string, error) {
	return Exec(fmt.Sprintf(format, args...))
}

func Exec(args ...string) (string, error) {
	list := append([]string{"/c"}, args...)
	switch runtime.GOOS {
	case "windows":
		result, err := exec.Command("cmd", list...).CombinedOutput()
		if err != nil {
			return "", err
		}
		result, err = str.GbkToUtf8(result)
		return string(result), err
	case "linux":
		list[0] = "-c"
		result, err := exec.Command("bash", list...).CombinedOutput()
		if err != nil {
			return "", err
		}
		return string(result), nil
	}
	return "", errors.New("未知操作系统:" + runtime.GOOS)
}

func Runf(format string, args ...interface{}) error {
	return Run(fmt.Sprintf(format, args...))
}

func Run(args ...string) error {
	list := append([]string{"/c"}, args...)
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", list...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	case "linux":
		return bash.Run(args...)
	}
	return errors.New("未知操作系统:" + runtime.GOOS)
}

func IO(w io.ReadWriter, args ...string) error {
	list := append([]string{"/c"}, args...)
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", list...)
		cmd.Stdout = w
		cmd.Stderr = w
		cmd.Stdin = w
		return cmd.Run()
	case "linux":
		return bash.IO(w, args...)
	}
	return errors.New("未知操作系统:" + runtime.GOOS)
}

// Stop 结束程序 "taskkill.exe", "/f", "/im", "edge.exe"
func Stop(name string) error {
	switch runtime.GOOS {
	case "windows":
		result, err := Exec("taskkill.exe", "/f", "/im", name)
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			return err
		} else if err == nil && !strings.Contains(result, "成功") {
			return errors.New(result)
		}
	case "linux":
		return systemctl.Stop(name)
	}
	return nil
}

// Start2 启动程序
// windows "cmd", "/c", "start ./xxx.exe"
func Start2(filename string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start "+filename).Start()
	case "linux":
		return systemctl.Restart(filename)
	}
	return nil
}
