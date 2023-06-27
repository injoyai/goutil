package bar

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Demo() {
	x := New().SetTotalSize(10)
	x.SetColor(color.FgBlue)
	go func() {
		for {
			time.Sleep(time.Second)
			x.Add(1)
		}
	}()
	x.Run()
}

func Download(url, filename string) error {
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

	return Copy(f, resp.Body, conv.Int64(resp.Header.Get("Content-Length")))
}

func Copy(w io.Writer, r io.Reader, total int64) error {
	buff := bufio.NewReader(r)
	b := New().SetTotalSize(float64(total))
	go b.Run()
	defer b.Done()
	for {
		buf := make([]byte, 1<<20)
		n, err := buff.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		b.Add(float64(n))
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
		if err == io.EOF {
			return nil
		}
	}
}

func New() *Bar {
	return &Bar{
		prefix:    "",
		suffix:    "",
		length:    40,
		nowSize:   0,
		totalSize: 1000,
		style:     ">",
		color:     nil, //color.New(color.Reset),
		c:         make(chan string, 1),
		done:      make(chan struct{}, 1),
		print:     func(s string) { fmt.Print(s) },
	}
}

type Bar struct {
	prefix    string        //前缀
	suffix    string        //后缀
	length    int           //总长度
	nowSize   float64       //当前完成数量
	totalSize float64       //最大数量
	style     string        //进度条风格
	color     *color.Color  //整体颜色
	c         chan string   //实时数据通道
	done      chan struct{} //结束信号
	print     func(string)  //打印
}

// SetPrint 设置打印函数
func (this *Bar) SetPrint(fn func(string)) *Bar {
	this.print = fn
	return this
}

// SetPrefix 设置前缀
func (this *Bar) SetPrefix(prefix string) *Bar {
	this.prefix = prefix
	return this
}

// SetWidth 设置进度条宽度
func (this *Bar) SetWidth(length int) *Bar {
	this.length = length
	return this
}

// SetTotal 设置进度任务数量
func (this *Bar) SetTotal(size float64) *Bar {
	return this.SetTotalSize(size)
}

// SetTotalSize 设置进度任务数量
func (this *Bar) SetTotalSize(size float64) *Bar {
	this.totalSize = size
	return this
}

// SetStyle 设置进度条风格
func (this *Bar) SetStyle(style string) *Bar {
	this.style = style
	return this
}

// SetColor 设置进度条颜色
func (this *Bar) SetColor(a color.Attribute) *Bar {
	this.color = color.New(a)
	return this
}

func (this *Bar) Done() {
	this.Add(this.totalSize - this.nowSize)
}

func (this *Bar) Add(n float64) {
	select {
	case <-this.done:
		return
	default:
		this.nowSize += n
		if this.nowSize >= this.totalSize && this.totalSize > 0 {
			this.nowSize = this.totalSize
			defer func() {
				select {
				case <-this.done:
				default:
					close(this.done)
				}
			}()
		}
		nowLength := int((this.nowSize / this.totalSize) * float64(this.length) / float64(len(this.style)))
		s := ""
		for i := 0; i < nowLength; i++ {
			s += this.style
		}
		this.c <- s
	}
}

func (this *Bar) Run() <-chan struct{} {
	this.Add(0)
	start := time.Now()
	for {
		select {
		case <-this.done:
			return this.done
		case s := <-this.c:
			width := strconv.Itoa(this.length)
			if this.color != nil {
				this.color.EnableColor()
				s = this.color.Sprint(s)
				width = strconv.Itoa(this.length + 9)
			}
			s = fmt.Sprintf("\r%s[%-"+width+"s] %0.1f%% %0.0f/%0.0f %0.0fs %s", this.prefix, s, this.nowSize*100/this.totalSize, this.nowSize, this.totalSize, time.Now().Sub(start).Seconds(), this.suffix)
			if this.print != nil {
				this.print(s)
			} else {
				fmt.Print(s)
			}
		}
	}
}
