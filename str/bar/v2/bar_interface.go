package bar

import (
	"fmt"
	"io"
	"time"
)

type Bar interface {
	fmt.Stringer
	io.Closer
	Add(n int64)                         //添加数量
	Set(current int64)                   //设置当前数量
	SetCurrent(current int64)            //设置当前数量
	SetTotal(total int64)                //设置总数量
	SetFormat(format func(b Bar) string) //设置样式
	SetPrefix(prefix string)             //设置前缀
	SetSuffix(suffix string)             //设置后缀
	SetWriter(w io.Writer)               //设置writer
	OnSet(f func())                      //设置事件
	OnFinal(f func(b Bar))               //完成事件

	Last() int64                          //最后数量
	Current() int64                       //当前数量
	Total() int64                         //总数量
	StartTime() time.Time                 //开始时间
	LastTime() time.Time                  //最后时间
	Flush() bool                          //刷入writer
	IntervalFlush(interval time.Duration) //间隔刷新
	Done() <-chan struct{}                //

	DownloadHTTP(source, filename string, proxy ...string) (int64, error) //通过http下载
	Copy(w io.Writer, r io.Reader) (int64, error)                         //复制
}

type Option func(b Bar)
type PlanOption func(p Plan)
type Format func(b Bar) string

type Plan interface {
	SetPrefix(s string) //设置前缀
	SetSuffix(s string) //设置后缀
	SetStyle(s byte)    //设置样式
	SetColor(c Color)   //设置颜色
	SetWidth(w int)     //设置宽度
}
