package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Interval 间隔时间
type Interval time.Duration

func (this Interval) String() string {
	return fmt.Sprintf("@every %s", time.Duration(this))
}

// NewIntervalSpec 新建间隔任务
func NewIntervalSpec(t time.Duration) string {
	return Interval(t).String()
}

// Date 按日志执行
type Date struct {
	Month  []int //月 1, 12
	Week   []int //周 0, 6
	Day    []int //天 1, 31
	Hour   []int //时 0, 23
	Minute []int //分 0, 59
	Second []int //秒 0, 59
}

func (this Date) spec(ints []int) string {
	if len(ints) > 0 {
		list := make([]string, len(ints))
		for i, v := range ints {
			list[i] = strconv.Itoa(v)
		}
		return strings.Join(list, ",")
	}
	return "*"
}

func (this Date) String() string {
	return strings.Join([]string{
		this.spec(this.Second),
		this.spec(this.Minute),
		this.spec(this.Hour),
		this.spec(this.Day),
		this.spec(this.Month),
		this.spec(this.Week),
	}, " ")
}
