package shell

import (
	"os/exec"
	"syscall"
)

func Start2(filename string) error {
	filename = `"" "` + filename + `"`
	cmd := exec.Command("cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "/c start " + filename}
	return cmd.Run()
}
