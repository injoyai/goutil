package shell

import (
	"github.com/injoyai/goutil/oss/linux/systemctl"
)

func Start(filename string) error {
	return systemctl.Restart(filename)
}
