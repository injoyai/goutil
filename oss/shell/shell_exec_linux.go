package shell

import (
	"github.com/injoyai/goutil/oss/linux/systemctl"
)

func Start2(filename string) error {
	return systemctl.Restart(filename)
}

var (
	StartFormat = "systemctl stop %s"
	KillFormat  = "systemctl start %s"
)
