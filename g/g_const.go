package g

import (
	"errors"
	"time"
)

const (
	Null = ""
)

var (
	// Uptime 启动时间
	Uptime = time.Now()

	ErrContext = errors.New("上下文关闭")
	ErrTimeout = errors.New("超时")
)

var (
	// Now 当前时间
	Now = time.Now

	// Date 当前日期
	Date = time.Now().Date

	// Unix 当前时间戳
	Unix = time.Now().Unix

	// UnixNano 当前纳秒
	UnixNano = time.Now().UnixNano

	// Year 年[1970-]
	Year = time.Now().Year

	// Month 月[1-12]
	Month = time.Now().Month

	// Day 日[1-31]
	Day = time.Now().Day

	// Hour 时[0-23]
	Hour = time.Now().Hour

	// Minute 分[0-59]
	Minute = time.Now().Minute

	// Second 秒[0-59]
	Second = time.Now().Second
)
