package bar

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"io"
	"os"
	"time"
)

var _ Interface = new(entity)

func New(total int64) Interface {
	return NewWithContext(context.Background(), total)
}

func NewWithContext(ctx context.Context, total int64) Interface {
	ctx, cancel := context.WithCancel(ctx)
	return &entity{
		width:  50,
		style:  '>',
		total:  total,
		ctx:    ctx,
		cancel: cancel,
		c:      make(chan int64),
	}
}

type entity struct {
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
}

func (this *entity) SetFormatter(f Formatter) Interface {
	this.format = f
	return this
}

func (this *entity) SetWidth(width int) Interface {
	this.width = width
	return this
}

func (this *entity) SetTotal(total int64) Interface {
	this.total = total
	return this
}

func (this *entity) SetStyle(style byte) Interface {
	this.style = style
	return this
}

func (this *entity) SetColor(a color.Attribute) Interface {
	this.color = color.New(a)
	return this
}

func (this *entity) Add(n int64) Interface {
	select {
	case <-this.ctx.Done():
	case this.c <- n:
	}
	return this
}

func (this *entity) Done() Interface {
	return this.Add(this.total)
}

func (this *entity) Close() error {
	this.cancel()
	return nil
}

func (this *entity) Run() <-chan struct{} {
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
			f := &Entity{
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
					currentNum, currentUnit := ToB(this.current)
					totalNum, totalUnit := ToB(this.total)
					return fmt.Sprintf("%0.1f%s/%0.1f%s", currentNum, currentUnit, totalNum, totalUnit)
				}),
				Speed: element(func() string {
					if val, ok := cache.Get("Speed"); ok {
						return val.(string)
					}
					f, unit := ToB(int64(spend))
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

			fmt.Print(s)

		}
	}
}

func (this *entity) init() {
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
		this.format = func(e *Entity) string {
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

func (this *entity) Copy(w io.Writer, r io.Reader) error {
	return this.CopyN(w, r, 4<<10)
}

func (this *entity) CopyN(w io.Writer, r io.Reader, num int64) error {
	defer this.Close()
	buff := bufio.NewReader(r)
	go this.Run()
	for {
		buf := make([]byte, num)
		n, err := buff.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		this.Add(int64(n))
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
		if err == io.EOF {
			return nil
		}
	}
}

var (
	defaultClient = http.NewClient()
)

func (this *entity) DownloadHTTP(source, filename string, proxy ...string) error {
	defaultClient.SetProxy(conv.GetDefaultString("", proxy...))
	resp := defaultClient.Get(source)
	if resp.Err() != nil {
		return resp.Err()
	}
	defer resp.Body.Close()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	total := conv.Int64(resp.GetHeader("Content-Length"))
	this.SetTotal(total)
	return this.Copy(f, resp.Body)
}
