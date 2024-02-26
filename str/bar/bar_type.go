package bar

import (
	"fmt"
	"github.com/fatih/color"
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
	Bar          *bar    // Element //进度条,例 [>>>   ]
	Rate         Element //进度百分比,例 58%
	RateSize     Element //进度数量,例 58/100
	RateSizeUnit Element //进度数量带单位,例 58B/100B
	Speed        Element //进度速度,例 13B/s
	SpeedUnit    Element //进度速度带单位,例 13MB/s
	Used         Element //已经耗时,例 2m20s
	Remain       Element //预计剩余时间 例 1m18s
}

// NewSizeUnit 自定义
func (this *Format) NewSizeUnit(size int64, decimal ...uint) string {
	f, unit := oss.SizeUnit(size)
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
		this.SpeedUnit,
	)
}

// M3u8 m3u8样式,进度是分片,字节是另外的
func (this *Format) M3u8(current int64) string {
	return fmt.Sprintf("\r%s  %s  %s  %s",
		this.Bar,
		this.RateSize,
		this.NewSizeUnit(current),
		this.Entity.SpeedUnit("m3u8", current, time.Millisecond*500),
	)
}

type Formatter func(e *Format) string

func WithDefault(e *Format) string {
	return e.Default()
}

func NewWithM3u8(current *int64) func(e *Format) string {
	return func(e *Format) string {
		return e.M3u8(*current)
	}
}

type bar struct {
	prefix, suffix string       //前缀后缀 例 []
	style          byte         //进度条风格 例 >
	color          *color.Color //整体颜色
	total          int64        //总数
	current        int64        //当前
	width          int          //宽度
}

func (this *bar) SetPrefix(prefix string) {
	this.prefix = prefix
}

func (this *bar) SetSuffix(suffix string) {
	this.suffix = suffix
}

func (this *bar) SetStyle(style byte) {
	this.style = style
}

func (this *bar) SetWidth(width int) {
	this.width = width
}

func (this *bar) SetColor(a color.Attribute) {
	this.color = color.New(a)
}

func (this *bar) String() string {
	rate := float64(this.current) / float64(this.total)
	nowWidth := ""
	for i := 0; i < int(float64(this.width)*rate); i++ {
		nowWidth += string(this.style)
	}
	barStr := fmt.Sprintf(fmt.Sprintf("%s%%-%ds%s", this.prefix, this.width, this.suffix), nowWidth)
	if this.color != nil {
		barStr = this.color.Sprint(barStr)
	}
	return barStr
}
