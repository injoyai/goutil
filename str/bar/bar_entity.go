package bar

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil"
	"io"
	"net/http"
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
	if this.c == nil {
		this.c = make(chan int64, 1)
	}
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
	if this.cancel != nil {
		this.cancel()
	}
	return nil
}

func (this *entity) Run() <-chan struct{} {
	this.init()
	start := time.Now()
	max := 0
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
					currentNum, currentUnit := this.toB(this.current)
					totalNum, totalUnit := this.toB(this.total)
					return fmt.Sprintf("%0.1f%s/%0.1f%s", currentNum, currentUnit, totalNum, totalUnit)
				}),
				Speed: element(func() string {
					// todo 算法待优化
					if spend <= 0 {
						spend = 0
					}
					f, unit := this.toB(int64(spend))
					return fmt.Sprintf("%0.1f%s/s", f, unit)
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
			if len(s) >= max {
				max = len(s)
			} else {
				s += fmt.Sprintf(fmt.Sprintf("%%-%ds", max-len(s)), " ")
			}

			fmt.Print(s)

		}
	}
}

func (this *entity) toB(n int64) (float64, string) {
	return goutil.ToB(n)
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
	this.Add(0)
}

func (this *entity) Copy(w io.Writer, r io.Reader) error {
	defer this.Close()
	buff := bufio.NewReader(r)
	go this.Run()
	for {
		buf := make([]byte, 2<<20)
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

func (this *entity) CopyN(w io.Writer, r io.Reader, num int64) error {
	buff := bufio.NewReader(r)
	go this.Run()
	for {
		buf := make([]byte, num)
		n, err := buff.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		this.Add(int64(n))
		num -= int64(n)
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
		if err == io.EOF || num <= 0 {
			return nil
		}
	}
}

func (this *entity) DownloadHTTP(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	total := conv.Int64(resp.Header.Get("Content-Length"))
	this.SetTotal(total)
	return this.Copy(f, resp.Body)
}
