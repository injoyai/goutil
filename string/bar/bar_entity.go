package bar

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"os"
	"time"
)

type Interface interface {

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
	SetFormatter(f Formatter)

	// SetPrefix 设置前缀,默认格式生效
	SetPrefix(prefix string)

	// SetSuffix 设置后缀,默认格式生效
	SetSuffix(suffix string)

	// SetWidth 设置宽度
	SetWidth(width int)

	// SetTotal 设置总数据大小
	SetTotal(total int64)

	// SetStyle 设置进度条风格
	SetStyle(style byte)

	// SetColor 设置颜色
	SetColor(color color.Attribute)

	// Add 添加数据
	Add(n int64)

	// Done 结束
	Done()

	// Run 运行
	Run() <-chan struct{}

	//衍生功能,

	// Copy 复制数据加入进度条
	Copy(w io.Writer, r io.Reader) error

	// CopyN 复制数据加入进度条
	CopyN(w io.Writer, r io.Reader, num int64) error

	// DownloadHTTP 下载http
	DownloadHTTP(url, filename string) error
}

var _ Interface = new(entity)

func New(total int64) Interface {
	return NewWithContext(context.Background(), total)
}

func NewWithContext(ctx context.Context, total int64) Interface {
	ctx, cancel := context.WithCancel(ctx)
	return &entity{
		width:  40,
		style:  '>',
		total:  total,
		ctx:    ctx,
		cancel: cancel,
	}
}

type entity struct {
	format  Formatter          //格式化
	prefix  string             //前缀
	suffix  string             //后缀
	width   int                //宽度
	current int64              //当前
	total   int64              //总
	style   byte               //进度条风格
	color   *color.Color       //整体颜色
	c       chan int64         //实时数据通道
	ctx     context.Context    //
	cancel  context.CancelFunc //
}

func (this *entity) SetFormatter(f Formatter) {
	this.format = f
}

func (this *entity) SetPrefix(prefix string) {
	this.prefix = prefix
}

func (this *entity) SetSuffix(suffix string) {
	this.suffix = suffix
}

func (this *entity) SetWidth(width int) {
	this.width = width
}

func (this *entity) SetTotal(total int64) {
	this.total = total
}

func (this *entity) SetStyle(style byte) {
	this.style = style
}

func (this *entity) SetColor(a color.Attribute) {
	this.color = color.New(a)
}

func (this *entity) Add(n int64) {
	if this.c == nil {
		this.c = make(chan int64, 1)
	}
	select {
	case <-this.ctx.Done():
	case this.c <- n:
	}
}

func (this *entity) Done() {
	this.Add(this.total)
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

			this.current += n
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
				Rate: element(func() string { return fmt.Sprintf("%0.1f%%", rate*100) }),
				Size: element(func() string { return fmt.Sprintf("%d/%d", this.current, this.total) }),
				Used: element(func() string { return fmt.Sprintf("%0.1fs", time.Now().Sub(start).Seconds()) }),
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
			return fmt.Sprintf("\r%s%s %s %s %s %s",
				this.prefix,
				e.Bar,
				e.Rate,
				e.Size,
				e.Used,
				this.suffix,
			)
		}
	}
	this.Add(0)
}

func (this *entity) Copy(w io.Writer, r io.Reader) error {
	buff := bufio.NewReader(r)
	go this.Run()
	for {
		buf := make([]byte, 1<<20)
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

type Element interface {
	String() string
}

type element func() string

func (this element) String() string { return this() }

type Entity struct {
	Bar    Element
	Rate   Element
	Size   Element
	Used   Element
	Remain Element
}

type Formatter func(e *Entity) string
