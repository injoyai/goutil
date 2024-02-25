package bar

import (
	"fmt"
	"github.com/injoyai/goutil/oss"
	"time"
)

type Element interface {
	String() string
}

type element func() string

func (this element) String() string { return this() }

type Format struct {
	Entity       *Bar    //实例
	Bar          Element //进度条,例 [>>>   ]
	Rate         Element //进度百分比,例 58%
	RateSize     Element //进度数量,例 58/100
	RateSizeUnit Element //进度数量带单位,例 58B/100B
	Speed        Element //进度速度,例 13B/s
	Used         Element //耗时,例 2m20s
	Remain       Element //预计剩余时间 例 1m18s
}

// NewSizeUnit 自定义
func (this *Format) NewSizeUnit(size int64, decimal ...uint) string {
	f, unit := oss.Size(size)
	if len(decimal) > 0 {
		return fmt.Sprintf(fmt.Sprintf("%%0.%df%%s", decimal[0]), f, unit)
	}
	return fmt.Sprintf("%0.1f%s", f, unit)
}

// Default 默认样式
func (this *Format) Default() string {
	return fmt.Sprintf("\r%s  %s  %s",
		this.Bar,
		this.RateSizeUnit,
		this.Speed,
	)
}

// M3u8 m3u8样式,进度是分片,字节是另外的
func (this *Format) M3u8(current int64) string {
	return fmt.Sprintf("\r%s  %s  %s  %s",
		this.Bar,
		this.RateSize,
		this.NewSizeUnit(current),
		this.Entity.Speed("m3u8", current, time.Millisecond*500),
	)
}

type Formatter func(e *Format) string

func WithDefault(e *Format) string {
	return e.Default()
}
