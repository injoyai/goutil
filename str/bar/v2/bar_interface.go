package bar

import (
	"fmt"
	"io"
	"time"
)

type Bar interface {
	Base
	Elements
}

type Base interface {
	fmt.Stringer
	io.Closer
	Add(n int64)           //添加数量
	Set(current int64)     //设置当前数量
	SetTotal(total int64)  //设置总数量
	SetWriter(w io.Writer) //设置writer
	OnFinal(f func(b Bar)) //完成事件

	Last() int64                          //最后数量
	Current() int64                       //当前数量
	Total() int64                         //总数量
	StartTime() time.Time                 //开始时间
	LastTime() time.Time                  //最后时间
	Flush() bool                          //刷入writer
	IntervalFlush(interval time.Duration) //间隔刷新

	DownloadHTTP(source, filename string, proxy ...string) (int64, error) //http下载
}

type Element = fmt.Stringer

type ElementFunc func() string

func (this ElementFunc) String() string { return this() }

type Elements interface {
	Plan(op ...PlanOption) Element //进度条,例 [>>>   ]
	Rate() Element                 //进度百分比,例 58%
	RateSize() Element             //进度数量,例 58/100
	RateSizeUnit() Element         //进度数量带单位,例 58B/100B
	Speed() Element                //进度速度,例 13B/s
	SpeedUnit() Element            //进度速度带单位,例 13MB/s
	SpeedAvg() Element             //进度平均速度,例 13B/s
	SpeedUnitAvg() Element         //进度平均速度带单位,例 13MB/s
	Used() Element                 //已经耗时,例 2m20s
	UsedSecond() Element           //已经耗时,例 600s
	Remain() Element               //预计剩余时间 例 1m18s
}

type Plan interface {
	SetPrefix(s string) //设置前缀
	SetSuffix(s string) //设置后缀
	SetStyle(s byte)    //设置样式
	SetColor(c Color)   //设置颜色
	SetWidth(w int)     //设置宽度
}

type PlanOption func(p Plan)
