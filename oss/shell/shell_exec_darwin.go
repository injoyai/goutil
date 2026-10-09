package shell

import (
	"os/exec"
)

// Start2 使用 open 启动程序
func Start2(filename string) error {
	return exec.Command("open", filename).Run()
}

var (
	StartFormat = "open %s"
	KillFormat  = "pkill -f %s"
)
