package bar

import (
	"bufio"
	"context"
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
		width:  50,
		style:  '>',
		total:  total,
		ctx:    ctx,
		cancel: cancel,
		c:      make(chan int64),
		writer: os.Stdout,
	}
}

type Bar struct {
	format      Formatter          //格式化
	width       int                //宽度
	current     int64              //当前
	currentTime time.Time          //当前时间
	total       int64              //总
	style       byte               //进度条风格
	color       *color.Color       //整体颜色
	c           chan int64         //实时数据通道
	ctx         context.Context    //
	cancel      context.CancelFunc //
	writer      io.Writer
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
*/
func (this *Bar) SetFormatter(f Formatter) *Bar {
	this.format = f
	return this
}

func (this *Bar) SetWidth(width int) *Bar {
	this.width = width
	return this
}

func (this *Bar) SetTotal(total int64) *Bar {
	this.total = total
	return this
}

func (this *Bar) SetStyle(style byte) *Bar {
	this.style = style
	return this
}

func (this *Bar) SetColor(a color.Attribute) *Bar {
	this.color = color.New(a)
	return this
}

func (this *Bar) SetWriter(w io.Writer) *Bar {
	this.writer = w
	return this
}

func (this *Bar) Add(n int64) *Bar {
	select {
	case <-this.ctx.Done():
	case this.c <- n:
	}
	return this
}

func (this *Bar) Done() *Bar {
	return this.Add(this.total)
}

func (this *Bar) Close() error {
	this.cancel()
	return nil
}

func (this *Bar) Run() <-chan struct{} {
	this.init()             //初始化
	go this.Add(0)          //触发进度条出现
	start := time.Now()     //开始时间
	maxLength := 0          //字符串最大长度
	cache := maps.NewSafe() //缓存,用于缓存最近的下载速度
	for {
		select {
		case <-this.ctx.Done():
			fmt.Println()
			return this.ctx.Done()
		case n := <-this.c:
			spend := float64(n) / time.Now().Sub(this.currentTime).Seconds()
			this.current += n
			this.currentTime = time.Now()
			if this.current >= this.total {
				this.current = this.total
				this.cancel()
			}

			//进度占比
			rate := float64(this.current) / float64(this.total)
			nowWidth := ""
			for i := 0; i < int(float64(this.width)*rate); i++ {
				nowWidth += string(this.style)
			}

			//元素
			f := &Format{
				Bar: element(func() string {
					bar := fmt.Sprintf(fmt.Sprintf("[%%-%ds]", this.width), nowWidth)
					if this.color != nil {
						bar = this.color.Sprint(bar)
					}
					return bar
				}),
				Rate: element(func() string {
					return fmt.Sprintf("%0.1f%%", rate*100)
				}),
				Size: element(func() string {
					return fmt.Sprintf("%d/%d", this.current, this.total)
				}),
				SizeUnit: element(func() string {
					currentNum, currentUnit := oss.Size(this.current)
					totalNum, totalUnit := oss.Size(this.total)
					return fmt.Sprintf("%0.1f%s/%0.1f%s", currentNum, currentUnit, totalNum, totalUnit)
				}),
				Speed: element(func() string {
					if val, ok := cache.Get("Speed"); ok {
						return val.(string)
					}
					f, unit := oss.Size(int64(spend))
					if f < 0 {
						f, unit = 0, "B"
					}
					s := fmt.Sprintf("%0.1f%s/s", f, unit)
					if f > 0 {
						cache.Set("Speed", s, time.Millisecond*500)
					}
					return s
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

func (this *Bar) init() {
	if this.width == 0 {
		this.width = 40
	}
	if this.style == 0 {
		this.style = '>'
	}
	if this.ctx == nil {
		this.ctx, this.cancel = context.WithCancel(context.Background())
	}
	if this.format == nil {
		this.format = func(e *Format) string {
			return fmt.Sprintf("\r%s  %s  %s",
				e.Bar,
				e.SizeUnit,
				e.Speed,
			)
		}
	}
	this.current = 0
	this.currentTime = time.Now()
}

func (this *Bar) Copy(w io.Writer, r io.Reader) (int, error) {
	return this.CopyN(w, r, 4<<10)
}

func (this *Bar) CopyN(w io.Writer, r io.Reader, bufSize int64) (int, error) {
	defer this.Close()
	buff := bufio.NewReader(r)
	go this.Run()
	total := 0
	for {
		buf := make([]byte, bufSize)
		n, err := buff.Read(buf)
		if err != nil && err != io.EOF {
			return total, err
		}
		total += n
		this.Add(int64(n))
		if _, err := w.Write(buf[:n]); err != nil {
			return total, err
		}
		if err == io.EOF {
			return total, nil
		}
	}
}

var (
	defaultClient = http.NewClient()
)

func (this *Bar) DownloadHTTP(source, filename string, proxy ...string) (int, error) {
	if err := defaultClient.SetProxy(conv.GetDefaultString("", proxy...)); err != nil {
		return 0, err
	}
	resp := defaultClient.Get(source)
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
