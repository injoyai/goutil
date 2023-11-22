package ntp

import (
	"github.com/injoyai/goutil/oss/linux/bash"
	"github.com/injoyai/goutil/str/regexps"
	"strings"
)

const (
	NTPAliyun       = "ntp.aliyun.com"
	NTPTencent      = "time1.cloud.tencent.com"
	DefaultNTP      = "ntp.ntsc.ac.cn"
	DefaultTimezone = "Asia/Shanghai"
)

// GetTimezoneList 获取可用的时区列表
func GetTimezoneList() ([]string, error) {
	result, err := bash.Exec("timedatectl list-timezones")
	return strings.Split(result, "\n"), err
}

// CurrentTimezone 当前时区
func CurrentTimezone() (string, error) {
	result, err := bash.Exec("timedatectl")
	if err != nil {
		return "", err
	}
	list := regexps.FindAll(`Time zone: [a-zA-Z/]+`, result)
	if len(list) > 0 {
		return strings.TrimLeft(list[0], "Time zone: "), nil
	}
	return "", nil
}

// SetTimezone 设置时区
func SetTimezone(timezone string) error {
	if len(timezone) == 0 {
		timezone = DefaultTimezone
	}
	_, err := bash.Exec("timedatectl set-timezone " + timezone)
	return err
}

// CurrentDate 获取当前日期
func CurrentDate() (string, error) {
	return bash.Exec("date")
}

// SyncDate 同步时间
func SyncDate(ntp string) error {
	if len(ntp) == 0 {
		//阿里ntp.aliyun.com 腾讯time1.cloud.tencent.com
		ntp = DefaultNTP
	}
	_, err := bash.Exec("ntpdate " + ntp)
	return err
}
