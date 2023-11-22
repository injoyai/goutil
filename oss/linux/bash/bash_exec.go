package bash

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Execf(format string, args ...interface{}) (string, error) {
	return Exec(fmt.Sprintf(format, args...))
}

func Exec(args ...string) (string, error) {
	list := append([]string{"-c"}, args...)
	result, err := exec.Command("bash", list...).CombinedOutput()
	if err != nil && len(result) > 0 {
		err = errors.New(string(result))
	}
	return string(result), err
}

func Runf(format string, args ...interface{}) error {
	return Run(fmt.Sprintf(format, args...))
}

func Run(args ...string) error {
	list := append([]string{"-c"}, args...)
	cmd := exec.Command("bash", list...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func IOf(w io.ReadWriter, format string, args ...interface{}) error {
	return IO(w, fmt.Sprintf(format, args...))
}

func IO(w io.ReadWriter, args ...string) error {
	list := append([]string{"-c"}, args...)
	cmd := exec.Command("bash", list...)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Stdin = w
	return cmd.Run()
}
