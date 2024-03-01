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

func (this *Bar) format() *Format {
	//进度占比
	rate := float64(this.current) / float64(this.total)
	return &Format{
		Entity: this,
		Bar: &bar{
			prefix:  "[",
			suffix:  "]",
			style:   '>',
			color:   nil,
			Total:   this.total,
			Current: this.current,
			width:   50,
		},
		Rate: element(func() string {
			return fmt.Sprintf("%0.1f%%", rate*100)
		}),
		RateSize: element(func() string {
			return fmt.Sprintf("%d/%d", this.current, this.total)
		}),
		RateSizeUnit: element(func() string {
			currentNum, currentUnit := oss.SizeUnit(this.current)
			totalNum, totalUnit := oss.SizeUnit(this.total)
			return fmt.Sprintf("%0.1f%s/%0.1f%s", currentNum, currentUnit, totalNum, totalUnit)
		}),
		Speed: element(func() string {
			return this.Speed("Speed", this.lastAdd, time.Millisecond*500)
		}),
		SpeedAvg: element(func() string {
			return this.SpeedAvg()
		}),
		SpeedUnit: element(func() string {
			return this.SpeedUnit("SpeedUnit", this.lastAdd, time.Millisecond*500)
		}),
		SpeedUnitAvg: element(func() string {
			return this.SpeedUnitAvg()
		}),
		Used: element(func() string {
			return fmt.Sprintf("%0.1fs", time.Now().Sub(this.start).Seconds())
		}),
		Remain: element(func() string {
			spend := time.Now().Sub(this.start)
			remain := "0s"
			if rate > 0 {
				remain = fmt.Sprintf("%0.1fs", time.Duration(float64(spend)/rate-float64(spend)).Seconds())
			}
			return remain
		}),
	}
}

type Format struct {
	Entity       *Bar    //实例
	Bar          *bar    // Element //进度条,例 [>>>   ]
	Rate         Element //进度百分比,例 58%
	RateSize     Element //进度数量,例 58/100
	RateSizeUnit Element //进度数量带单位,例 58B/100B
	Speed        Element //进度速度,例 13B/s
	SpeedAvg     Element //进度平均速度,例 13B/s
	SpeedUnit    Element //进度速度带单位,例 13MB/s
	SpeedUnitAvg Element //进度平均速度带单位,例 13MB/s
	Used         Element //已经耗时,例 2m20s
	Remain       Element //预计剩余时间 例 1m18s
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
func (this *Format) M3u8(once, current int64) string {
	return fmt.Sprintf("\r%s  %s  %s  %s",
		this.Bar,
		this.RateSize,
		this.Entity.SizeUnit(current),
		this.Entity.SpeedUnit("m3u8", once, time.Millisecond*500),
	)
}

type Formatter func(e *Format) string

func WithDefault(e *Format) string {
	return e.Default()
}

func NewWithM3u8(once, current *int64) func(e *Format) string {
	return func(e *Format) string {
		return e.M3u8(*once, *current)
	}
}

type bar struct {
	prefix, suffix string       //前缀后缀 例 []
	style          byte         //进度条风格 例 >
	color          *color.Color //整体颜色
	Total          int64        //总数
	Current        int64        //当前
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
	rate := float64(this.Current) / float64(this.Total)
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
