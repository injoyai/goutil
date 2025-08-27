package bar

import (
	"bufio"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"io"
	"os"
	"strings"
	"time"
)

// WithCurrent 设置当前数量
func WithCurrent(current int64) Option {
	return func(b Bar) {
		b.SetCurrent(current)
	}
}

// WithTotal 设置总数量
func WithTotal(total int64) Option {
	return func(b Bar) {
		b.SetTotal(total)
	}
}

// WithFormat 设置样式
func WithFormat(fs ...Format) Option {
	switch len(fs) {
	case 0:
		return func(b Bar) {}
	case 1:
		return func(b Bar) { b.SetFormat(fs[0]) }
	default:
		ls := make([]string, len(fs))
		return func(b Bar) {
			b.SetFormat(func(b Bar) string {
				for i, v := range fs {
					ls[i] = v(b)
				}
				return strings.Join(ls, "  ")
			})
		}

	}
}

// WithFormatDefault 设置默认样式,不带单位
func WithFormatDefault(op ...PlanOption) Option {
	return WithFormat(
		WithPlan(op...),
		WithRateSize(),
		WithSpeed(),
	)
}

// WithFormatDefaultUnit 设置默认样式,带单位
func WithFormatDefaultUnit(op ...PlanOption) Option {
	return WithFormat(
		WithPlan(op...),
		WithRateSizeUnit(),
		WithSpeedUnit(),
	)
}

// WithPrefix 设置前缀
func WithPrefix(prefix string) Option {
	return func(b Bar) {
		b.SetPrefix(prefix)
	}
}

// WithSuffix 设置后缀
func WithSuffix(suffix string) Option {
	return func(b Bar) {
		b.SetSuffix(suffix)
	}
}

// WithWriter 设置writer
func WithWriter(writer io.Writer) Option {
	return func(b Bar) {
		b.SetWriter(writer)
	}
}

// WithAutoFlush 设置后自动刷新
func WithAutoFlush() Option {
	return func(b Bar) {
		b.OnSet(func() {
			b.Flush()
		})
	}
}

// WithIntervalFlush 设置定时刷新
func WithIntervalFlush(interval time.Duration) Option {
	return func(b Bar) {
		b.IntervalFlush(interval)
	}
}

// WithFlush 刷入writer
func WithFlush() Option {
	return func(b Bar) {
		b.Flush()
	}
}

func New(op ...Option) Bar {
	b := &base{
		current:   0,
		total:     0,
		writer:    os.Stdout,
		Closer:    safe.NewCloser(),
		startTime: time.Now(),
	}
	b.SetCloseFunc(func(err error) error {
		if b.writer != nil {
			b.writer.Write([]byte("\n"))
		}
		return nil
	})
	WithFormatDefault()(b)
	for _, v := range op {
		v(b)
	}
	return b
}

type base struct {
	current      int64       //当前数量
	total        int64       //总数量
	prefix       string      //前缀
	suffix       string      //后缀
	format       Format      //格式化
	writer       io.Writer   //输出
	*safe.Closer             //closer
	onSet        func()      //设置事件
	onFinal      func(b Bar) //完成事件

	startTime time.Time //开始时间
	last      int64     //最后一次增加的值
	lastTime  time.Time //最后一次时间
}

func (this *base) Add(n int64) {
	this.Set(this.Current() + n)
}

func (this *base) Set(current int64) {
	this.SetCurrent(current)
}

func (this *base) SetCurrent(current int64) {
	if current > this.total {
		current = this.total
	}
	this.last = current - this.current
	this.lastTime = time.Now()
	this.current = current
	if this.onSet != nil {
		this.onSet()
	}
}

func (this *base) SetTotal(total int64) {
	this.total = total
}

func (this *base) SetFormat(format func(b Bar) string) {
	if format != nil {
		this.format = format
	}
}

func (this *base) SetPrefix(prefix string) {
	this.prefix = prefix
}

func (this *base) SetSuffix(suffix string) {
	this.suffix = suffix
}

func (this *base) SetWriter(w io.Writer) {
	this.writer = w
}

func (this *base) OnSet(f func()) {
	this.onSet = f
}

func (this *base) OnFinal(f func(b Bar)) {
	this.onFinal = f
}

/*



 */

func (this *base) Last() int64 {
	return this.last
}

func (this *base) Current() int64 {
	return this.current
}

func (this *base) Total() int64 {
	return this.total
}

func (this *base) StartTime() time.Time {
	return this.startTime
}

func (this *base) LastTime() time.Time {
	return this.lastTime
}

func (this *base) Flush() (closed bool) {
	if this.Closed() {
		return true
	}
	if this.writer == nil {
		return false
	}
	s := this.String()
	if s == "" || s[0] != '\r' {
		s = "\r\033[K" + s
	}
	this.writer.Write([]byte(s))
	if this.current >= this.total {
		if this.onFinal != nil {
			this.onFinal(this)
		}
		this.Close()
	}
	return this.Closed()
}

func (this *base) IntervalFlush(interval time.Duration) {
	for {
		select {
		case <-this.Done():
			return
		case <-time.After(interval):
			this.Flush()
		}
	}
}

func (this *base) String() string {
	return this.prefix + this.format(this) + this.suffix
}

/*



 */

/*



 */

var (
	// DefaultClient 默认客户端,下载大文件的时候需要设置长的超时时间
	DefaultClient = http.NewClient().SetTimeout(0)
)

func (this *base) DownloadHTTP(source, filename string, proxy ...string) (int64, error) {
	if err := DefaultClient.SetProxy(conv.Default("", proxy...)); err != nil {
		return 0, err
	}
	defer this.Close()
	return DefaultClient.GetToFileWithPlan(source, filename, func(p *http.Plan) {
		this.SetTotal(p.Total)
		this.Set(p.Current)
		this.Flush()
	})
}

func (this *base) Copy(w io.Writer, r io.Reader) (int64, error) {
	return this.CopyN(w, r, 4<<10)
}

func (this *base) CopyN(w io.Writer, r io.Reader, bufSize int64) (int64, error) {
	buff := bufio.NewReader(r)
	reader := &Read{
		Reader: buff,
		Bar:    this,
	}
	return io.CopyN(w, reader, bufSize)
	//total := int64(0)
	//buf := make([]byte, bufSize)
	//for {
	//	n, err := buff.Read(buf)
	//	if err != nil && err != io.EOF {
	//		return total, err
	//	}
	//	total += int64(n)
	//	this.Add(int64(n))
	//	this.Flush()
	//	if _, err := w.Write(buf[:n]); err != nil {
	//		return total, err
	//	}
	//	if err == io.EOF {
	//		return total, nil
	//	}
	//}
}

type Read struct {
	io.Reader
	Bar
}

func (this *Read) Read(p []byte) (n int, err error) {
	n, err = this.Reader.Read(p)
	this.Bar.Add(int64(n))
	this.Bar.Flush()
	return
}
