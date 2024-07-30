package ntp

import (
	"fmt"
	"github.com/beevik/ntp"
	"os/exec"
	"time"
)

func Get(s string) (*ntp.Response, error) {
	return ntp.Query(s)
}

// Sync 从NTP服务器上同步时间,并通过date命令设置时间
func Sync(s ...string) error {
	ser := NTPAliyun
	if len(s) > 0 && len(s[0]) > 0 {
		ser = s[0]
	}
	resp, err := ntp.Query(ser)
	if err != nil {
		return err
	}
	t := time.Now().Add(resp.ClockOffset)
	_, err = exec.Command("sh", "-c", fmt.Sprintf(`date -s "%s"`, t.Format(time.DateTime))).Output()
	return err
}
