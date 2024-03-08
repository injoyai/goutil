package shell

import (
	"os/exec"
	"syscall"
)

func Start(filename string) error {
	filename = `"" "` + filename + `"`
	cmd := exec.Command("cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "/c start " + filename}
	return cmd.Run()
}
