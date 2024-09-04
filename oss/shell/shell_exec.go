package shell

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/injoyai/goutil/str"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

var (
	Bash       = &Shell{_bash{}}
	SH         = &Shell{_sh{}}
	CMD        = &Shell{_cmd{}}
	Default    *Shell
	defaultUse Use
)

func init() {
	if runtime.GOOS == "windows" {
		defaultUse = _cmd{}
	} else {
		defaultUse = _bash{}
	}
	Default = &Shell{defaultUse}
}

func Execf(format string, args ...interface{}) (*Result, error) {
	return Default.Execf(format, args...)
}

func Exec(args ...string) (*Result, error) {
	return Default.Exec(args...)
}

func Run(args ...string) error {
	return Default.Run(args...)
}

func Output(w io.ReadWriter, args ...string) error {
	return Default.Output(w, args...)
}

// Stop 结束程序 "taskkill.exe", "/f", "/im", "edge.exe"
func Stop(name string) error {
	return Default.Stop(name)
}

// Start 启动程序
// windows "cmd", "/c", "start ./xxx.exe"
func Start(filename string) error {
	return Default.Start(filename)
}

/*



 */

type _cmd struct{}

func (_cmd) Prefix() [2]string { return [2]string{"cmd", "/c"} }

func (_cmd) Decode(p []byte) ([]byte, error) { return str.GbkToUtf8(p) }

type _bash struct{}

func (_bash) Prefix() [2]string { return [2]string{"bash", "-c"} }

func (_bash) Decode(p []byte) ([]byte, error) { return p, nil }

type _sh struct{}

func (_sh) Prefix() [2]string { return [2]string{"sh", "-c"} }

func (_sh) Decode(p []byte) ([]byte, error) { return p, nil }

type Use interface {
	Prefix() [2]string
	Decode(p []byte) ([]byte, error)
}

type Shell struct {
	Use
}

func (this *Shell) Execf(format string, args ...interface{}) (*Result, error) {
	cmdline := fmt.Sprintf(format, args...)
	return this.Exec(cmdline)
}

func (this *Shell) Exec(args ...string) (*Result, error) {
	pre := this.Prefix()
	list := append(pre[1:], args...)
	cmd := exec.Command(pre[0], list...)
	result := &Result{
		buf:    bytes.NewBuffer(nil),
		decode: this.Decode,
	}
	cmd.Stdout = result.buf
	cmd.Stderr = result.buf
	return result, cmd.Run()
}

func (this *Shell) Run(args ...string) error {
	pre := this.Prefix()
	list := append(pre[1:], args...)
	cmd := exec.Command(pre[0], list...)
	return cmd.Run()
}

func (this *Shell) Output(w io.Writer, args ...string) error {
	pre := this.Prefix()
	list := append(pre[1:], args...)
	cmd := exec.Command(pre[0], list...)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

func (this *Shell) Timeout(t time.Duration, args ...string) (*Result, error) {
	pre := this.Prefix()
	list := append(pre[1:], args...)
	cmd := exec.Command(pre[0], list...)
	result := &Result{
		buf:    bytes.NewBuffer(nil),
		decode: this.Decode,
	}
	cmd.Stdout = result.buf
	cmd.Stderr = result.buf

	var err error
	timer := time.NewTimer(t)
	defer timer.Stop()
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err = <-done:

	case <-timer.C:
		if err = cmd.Process.Kill(); err != nil {
			err = fmt.Errorf("执行超时,关闭命令失败: %s", err.Error())
		} else {
			// wait for the command to return after killing it
			<-done
			err = errors.New("执行超时")
		}
	}

	return result, err
}

// Start  "cmd", "/c", "start ./xxx.exe"
func (this *Shell) Start(filename string) error {
	_, err := this.Execf(StartFormat, filename)
	return err
}

// Stop 结束程序 "taskkill.exe", "/f", "/im", "edge.exe"
func (this *Shell) Stop(name string) error {
	result, err := this.Execf(KillFormat, name)
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		return err
	} else if err == nil && !strings.Contains(result.String(), "成功") {
		return errors.New(result.String())
	}
	return nil
}

type Result struct {
	buf    *bytes.Buffer
	str    *string
	decode func(p []byte) ([]byte, error)
}

func (this *Result) String() string {
	if this.str == nil {
		if this.decode != nil {
			bs, err := this.decode(this.buf.Bytes())
			if err == nil {
				this.str = (*string)(unsafe.Pointer(&bs))
				return *this.str
			}
		}
		this.str = (*string)(unsafe.Pointer(this.buf))
	}
	return *this.str
}
