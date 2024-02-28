package bar

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/goutil/oss"
	"io"
	"os"
	"time"
)

func New(total int64) *Bar {
	return NewWithContext(context.Background(), total)
}

func NewWithContext(ctx context.Context, total int64) *Bar {
	ctx, cancel := context.WithCancel(ctx)
	return &Bar{
		total:      total,
		ctx:        ctx,
		cancel:     cancel,
		c:          make(chan int64),
		writer:     os.Stdout,
		cacheTime:  maps.NewSafe(),
		cacheSpeed: maps.NewSafe(),
	}
}

type Bar struct {
	format Formatter //格式化
	option []func(format *Format)

	current    int64              //当前数量
	total      int64              //总数量
	c          chan int64         //实时数据通道
	ctx        context.Context    //
	cancel     context.CancelFunc //
	writer     io.Writer
	cacheTime  *maps.Safe
	cacheSpeed *maps.Safe
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
	this.format = f
	return this
}

// SetTotal 设置进度条的总数量,不影响宽度
func (this *Bar) SetTotal(total int64) *Bar {
	this.total = total
	return this
}

// AddOption 添加option
func (this *Bar) AddOption(option ...func(f *Format)) *Bar {
	this.option = append(this.option, option...)
	return this
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

// SetWriter 设置输出,默认到控制台
func (this *Bar) SetWriter(w io.Writer) *Bar {
	this.writer = w
	return this
}

// Write 实现io.Writer,取字节的长度
func (this *Bar) Write(p []byte) (int, error) {
	length := len(p)
	this.Add(int64(length))
	return length, nil
}

// Add 增加当前数量
func (this *Bar) Add(n int64) *Bar {
	select {
	case <-this.ctx.Done():
	case this.c <- n:
	}
	return this
}

// Done 完成执行,一次性添加剩余的数量
func (this *Bar) Done() *Bar {
	return this.Add(this.total - this.current)
}

// Close 结束执行,不添加剩余数量
func (this *Bar) Close() error {
	this.cancel()
	return nil
}

// Speed 计算速度,不带单位
func (this *Bar) Speed(key string, size int64, expiration time.Duration) string {
	return this.speed(key, size, expiration, func(size float64) string {
		return fmt.Sprintf("%0.1f/s", size)
	})
}

// SpeedUnit 速度,带单位 1024 -> 1KB
func (this *Bar) SpeedUnit(key string, size int64, expiration time.Duration) string {
	return this.speed(key, size, expiration, func(size float64) string {
		f, unit := oss.SizeUnit(int64(size))
		return fmt.Sprintf("%0.1f%s/s", f, unit)
	})
}

func (this *Bar) Run() error {
	defer fmt.Println()
	this.current = 0
	go this.Add(0)      //触发进度条出现
	start := time.Now() //开始时间
	maxLength := 0      //字符串最大长度
	for {
		select {
		case <-this.ctx.Done():
			return errors.New("上下文关闭")

		case n := <-this.c:

			//当前进度
			this.current += n
			if this.current >= this.total {
				this.current = this.total
				this.cancel()
			}

			//进度占比
			rate := float64(this.current) / float64(this.total)

			//元素
			f := &Format{
				Entity: this,
				Bar: &bar{
					prefix:  "[",
					suffix:  "]",
					style:   '>',
					color:   nil,
					total:   this.total,
					current: this.current,
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
					return this.Speed("Speed", n, time.Millisecond*500)
				}),
				SpeedUnit: element(func() string {
					return this.SpeedUnit("SpeedUnit", n, time.Millisecond*500)
				}),
				Used: element(func() string {
					return fmt.Sprintf("%0.1fs", time.Now().Sub(start).Seconds())
				}),
				Remain: element(func() string {
					spend := time.Now().Sub(start)
					remain := "0s"
					if rate > 0 {
						remain = fmt.Sprintf("%0.1fs", time.Duration(float64(spend)/rate-float64(spend)).Seconds())
					}
					return remain
				}),
			}

			for _, v := range this.option {
				v(f)
			}
			//自定义格式化输出
			if this.format == nil {
				this.format = WithDefault
			}
			s := this.format(f)
			if len(s) >= maxLength {
				maxLength = len(s)
			} else {
				s += fmt.Sprintf(fmt.Sprintf("%%-%ds", maxLength-len(s)), " ")
			}

			if this.writer != nil {
				this.writer.Write([]byte(s))
			}

		}
	}
}

func (this *Bar) Copy(w io.Writer, r io.Reader) (int64, error) {
	return this.CopyN(w, r, 4<<10)
}

func (this *Bar) CopyN(w io.Writer, r io.Reader, bufSize int64) (int64, error) {
	defer this.Close()
	buff := bufio.NewReader(r)
	go this.Run()
	total := int64(0)
	buf := make([]byte, bufSize)
	for {
		n, err := buff.Read(buf)
		if err != nil && err != io.EOF {
			return total, err
		}
		total += int64(n)
		this.Add(int64(n))
		if _, err := w.Write(buf[:n]); err != nil {
			return total, err
		}
		if err == io.EOF {
			return total, nil
		}
	}
}

func (this *Bar) speed(key string, size int64, expiration time.Duration, fn func(float64) string) string {

	//最后的数据时间
	lastTime, _ := this.cacheTime.GetOrSetByHandler(key, func() (interface{}, error) {
		return time.Time{}, nil
	})

	//记录这次时间,用于下次计算时间差
	now := time.Now()
	this.cacheTime.Set(key, now)

	//尝试从缓存获取速度,存在则直接返回,由expiration控制
	if val, ok := this.cacheSpeed.Get(key); ok {
		return val.(string)
	}

	//计算速度
	spendSize := float64(size) / now.Sub(lastTime.(time.Time)).Seconds()
	s := fn(spendSize)
	this.cacheSpeed.Set(key, s, expiration)
	return s
}

var (
	// DefaultClient 默认客户端,下载大文件的时候需要设置长的超时时间
	DefaultClient = http.NewClient().SetTimeout(0)
)

func (this *Bar) DownloadHTTP(source, filename string, proxy ...string) (int64, error) {
	if err := DefaultClient.SetProxy(conv.GetDefaultString("", proxy...)); err != nil {
		return 0, err
	}
	resp := DefaultClient.Get(source)
	if resp.Err() != nil {
		return 0, resp.Err()
	}
	defer resp.Body.Close()
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	total := conv.Int64(resp.GetHeader("Content-Length"))
	this.SetTotal(total)
	return this.Copy(f, resp.Body)
}
