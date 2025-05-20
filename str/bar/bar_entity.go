package bar

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/goutil/oss"
	"io"
	"os"
	"time"
)

func New(total int64) *Bar {
	return &Bar{
		total:      total,
		writer:     os.Stdout,
		cacheSpeed: maps.NewSafe(),
		start:      time.Now(),
		Closer:     safe.NewCloser(),
	}
}

type Bar struct {
	option       []func(format *Format) //选项
	formatter    Formatter              //格式化
	current      int64                  //当前数量
	total        int64                  //总数量
	writer       io.Writer              //输出
	*safe.Closer                        //

	cacheSpeed *maps.Safe //用于计算实时速度
	start      time.Time  //开始时间
	maxLength  int        //字符最大长度
	lastAdd    int64      //最后一次增加的数量
}

/*
SetFormatter 自定义格式

	func(e *Entity) string {
		return fmt.Sprintf("\r%s%s %s %s %s %s",
			this.prefix,
			e.Bar,
			e.Rate,
			e.Size,
			e.Remain,
			this.suffix,
		)
	}

这里可以也直接打印,需要自行处理上次数据遗留问题
*/
func (this *Bar) SetFormatter(f Formatter) *Bar {
	this.formatter = f
	return this
}

// SetTotal 设置进度条的总数量,不影响宽度
func (this *Bar) SetTotal(total int64) *Bar {
	this.total = total
	return this
}

// Set 同SetCurrent
func (this *Bar) Set(current int64) *Bar {
	return this.SetCurrent(current)
}

// SetCurrent 设置当前数量
func (this *Bar) SetCurrent(current int64) *Bar {
	if current > this.total {
		current = this.total
	}
	this.lastAdd = current - this.current
	this.current = current
	return this
}

// SetCurrentFlush 设置当前数值,并刷新
func (this *Bar) SetCurrentFlush(current int64) bool {
	return this.SetCurrent(current).Flush()
}

// Add 同AddCurrent
func (this *Bar) Add(size int64) *Bar {
	return this.AddCurrent(size)
}

// AddCurrent 增加当前数量
func (this *Bar) AddCurrent(size int64) *Bar {
	return this.SetCurrent(this.current + size)
}

// AddCurrentFlush 增加数值,并刷新
func (this *Bar) AddCurrentFlush(size int64) bool {
	return this.AddCurrent(size).Flush()
}

// SetWriter 设置输出,默认到控制台
func (this *Bar) SetWriter(w io.Writer) *Bar {
	this.writer = w
	return this
}

// Write 实现io.Writer,取字节的长度
func (this *Bar) Write(p []byte) (int, error) {
	length := len(p)
	this.AddCurrentFlush(int64(length))
	return length, nil
}

// Flush 刷新输出,返回是否关闭(结束)
func (this *Bar) Flush() (closed bool) {
	if this.Closed() {
		return true
	}
	if this.writer == nil {
		return false
	}
	s := this.String()
	if !(len(s) > 0 && s[0] == '\r') {
		s = "\r" + s
	}
	this.writer.Write([]byte(s))
	if this.current >= this.total {
		this.Close()
	}
	return this.Closed()
}

// Final 是否最后
func (this *Bar) Final() bool {
	return this.current >= this.total
}

// Close 关闭,输出换行符,后续不会继续输出
func (this *Bar) Close() error {
	if this.writer != nil {
		this.writer.Write([]byte("\n"))
	}
	return this.Closer.Close()
}

func (this *Bar) String() string {
	f := this.format()
	for _, v := range this.option {
		v(f)
	}
	//自定义格式化输出
	if this.formatter == nil {
		this.formatter = WithDefault
	}
	s := this.formatter(f)
	if len(s) >= this.maxLength {
		this.maxLength = len(s)
	} else {
		s += fmt.Sprintf(fmt.Sprintf("%%-%ds", this.maxLength-len(s)), " ")
	}
	return s
}

/*



 */

// AddOption 添加option
func (this *Bar) AddOption(option ...func(f *Format)) *Bar {
	this.option = append(this.option, option...)
	return this
}

// OnFinal 执行结束时
func (this *Bar) OnFinal(fn func(f *Format)) *Bar {
	return this.AddOption(func(f *Format) {
		if f.Entity.Final() {
			fn(f)
		}
	})
}

// SetWidth 设置进度条的宽度
func (this *Bar) SetWidth(width int) *Bar {
	return this.AddOption(func(format *Format) {
		format.Bar.SetWidth(width)
	})
}

// SetStyle 设置进度条风格,[>>>   ] [###   ] 等等
func (this *Bar) SetStyle(style byte) *Bar {
	return this.AddOption(func(format *Format) {
		format.Bar.SetStyle(style)
	})
}

// SetColor 设置进度条颜色
func (this *Bar) SetColor(a color.Attribute) *Bar {
	return this.AddOption(func(format *Format) {
		format.Bar.SetColor(a)
	})
}

/*



 */

// SizeUnit 字节单位1024 -> 1KB
func (this *Bar) SizeUnit(size int64, decimal ...uint) string {
	f, unit := oss.SizeUnit(size)
	if len(decimal) > 0 {
		return fmt.Sprintf(fmt.Sprintf("%%0.%df%%s", decimal[0]), f, unit)
	}
	return fmt.Sprintf("%0.1f%s", f, unit)
}

// Speed 实时速度,不带单位
func (this *Bar) Speed(key string, size int64, expiration time.Duration) string {
	return this.speed(key, size, expiration, func(size float64) string {
		return fmt.Sprintf("%0.1f/s", size)
	})
}

// SpeedAvg 平均速度,不带单位
func (this *Bar) SpeedAvg() string {
	speedSize := float64(this.current) / time.Since(this.start).Seconds()
	return fmt.Sprintf("%0.1f/s", speedSize)
}

// SpeedUnit 实时速度,带单位 1024/s -> 1KB/s
func (this *Bar) SpeedUnit(key string, size int64, expiration time.Duration) string {
	return this.speed(key, size, expiration, func(size float64) string {
		f, unit := oss.SizeUnit(int64(size))
		return fmt.Sprintf("%0.1f%s/s", f, unit)
	})
}

// SpeedUnitAvg 平均速度,带单位 1024/s -> 1KB/s
func (this *Bar) SpeedUnitAvg() string {
	speedSize := float64(this.current) / time.Since(this.start).Seconds()
	f, unit := oss.SizeUnit(int64(speedSize))
	return fmt.Sprintf("%0.1f%s/s", f, unit)
}

func (this *Bar) Copy(w io.Writer, r io.Reader) (int64, error) {
	return this.CopyN(w, r, 4<<10)
}

func (this *Bar) CopyN(w io.Writer, r io.Reader, bufSize int64) (int64, error) {
	buff := bufio.NewReader(r)
	total := int64(0)
	buf := make([]byte, bufSize)
	for {
		n, err := buff.Read(buf)
		if err != nil && err != io.EOF {
			return total, err
		}
		total += int64(n)
		this.Add(int64(n)).Flush()
		if _, err := w.Write(buf[:n]); err != nil {
			return total, err
		}
		if err == io.EOF {
			return total, nil
		}
	}
}

var (
	// DefaultClient 默认客户端,下载大文件的时候需要设置长的超时时间
	DefaultClient = http.NewClient().SetTimeout(0)
)

func (this *Bar) DownloadHTTP(source, filename string, proxy ...string) (int64, error) {
	if err := DefaultClient.SetProxy(conv.Default[string]("", proxy...)); err != nil {
		return 0, err
	}
	defer this.Close()
	return DefaultClient.GetToFileWithPlan(source, filename, func(p *http.Plan) {
		this.SetTotal(p.Total)
		this.Set(p.Current).Flush()
	})
}

func (this *Bar) speed(key string, size int64, expiration time.Duration, fn func(float64) string) string {

	timeKey := "time_" + key
	cacheKey := "speed_" + key
	//最后的数据时间
	lastTime, _ := this.cacheSpeed.GetOrSetByHandler(timeKey, func() (any, error) {
		return time.Time{}, nil
	})

	//记录这次时间,用于下次计算时间差
	now := time.Now()
	this.cacheSpeed.Set(timeKey, now)

	//尝试从缓存获取速度,存在则直接返回,由expiration控制
	if val, ok := this.cacheSpeed.Get(cacheKey); ok {
		return val.(string)
	}

	//计算速度
	size = conv.Select[int64](size >= 0, size, 0)
	spendSize := float64(size) / now.Sub(lastTime.(time.Time)).Seconds()
	s := fn(spendSize)
	this.cacheSpeed.Set(cacheKey, s, expiration)
	return s
}
