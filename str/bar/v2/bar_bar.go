package bar

import (
	"bufio"
	"fmt"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/goutil/oss"
	"io"
	"os"
	"time"
)

func WithTotal(total int64) Option {
	return func(b Base) {
		b.SetTotal(total)
	}
}

func WithFormat(format func(b Bar) string) Option {
	return func(b Base) {
		b.SetFormat(format)
	}
}

func WithDefaultFormat(op ...PlanOption) Option {
	return WithFormat(func(b Bar) string {
		return fmt.Sprintf("\r%s  %s  %s",
			b.Plan(op...),
			b.RateSizeUnit(),
			b.SpeedUnit(),
		)
	})
}

func WithWriter(writer io.Writer) Option {
	return func(b Base) {
		b.SetWriter(writer)
	}
}

type Option func(b Base)

func New(op ...Option) Bar {
	b := &base{
		current: 0,
		total:   0,
		format:  DefaultFormat,
		writer:  os.Stdout,
		Closer:  safe.NewCloser(),

		cache:     maps.NewSafe(),
		startTime: time.Now(),
	}
	b.SetCloseFunc(func(err error) error {
		if b.writer != nil {
			b.writer.Write([]byte("\n"))
		}
		return nil
	})
	for _, v := range op {
		v(b)
	}
	return b
}

type base struct {
	current      int64              //当前数量
	total        int64              //总数量
	format       func(b Bar) string //格式化
	writer       io.Writer          //输出
	*safe.Closer                    //

	cache     *maps.Safe //
	startTime time.Time  //开始时间
	last      int64      //最后一次增加的值
	lastTime  time.Time  //最后一次时间

	onFinal func(b Bar)
}

func (this *base) Add(n int64) {
	this.Set(this.Current() + n)
}

func (this *base) Set(current int64) {
	if current > this.total {
		current = this.total
	}
	this.last = current - this.current
	this.lastTime = time.Now()
	this.current = current
}

func (this *base) SetTotal(total int64) {
	this.total = total
}

func (this *base) SetFormat(format func(b Bar) string) {
	this.format = format
}

func (this *base) SetWriter(w io.Writer) {
	this.writer = w
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
		s = "\r" + s
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

var DefaultFormat = func(b Bar) string {
	return fmt.Sprintf("\r%s  %s  %s",
		b.Plan(),
		b.RateSizeUnit(),
		b.SpeedUnit(),
	)
}

func (this *base) String() string {
	if this.format == nil {
		this.format = DefaultFormat
	}
	return this.format(this)
}

/*



 */

func (this *base) Plan(op ...PlanOption) Element {
	b := &plan{
		prefix:  "[",
		suffix:  "]",
		style:   '>',
		color:   nil,
		width:   50,
		current: this.current,
		total:   this.total,
	}
	for _, v := range op {
		v(b)
	}
	return b
}

func (this *base) Rate() Element {
	return ElementFunc(func() string {
		return fmt.Sprintf("%0.1f%%", float64(this.current)*100/float64(this.total))
	})
}

func (this *base) RateSize() Element {
	return ElementFunc(func() string {
		return fmt.Sprintf("%d/%d", this.current, this.total)
	})
}

func (this *base) RateSizeUnit() Element {
	return ElementFunc(func() string {
		currentNum, currentUnit := oss.SizeUnit(this.current)
		totalNum, totalUnit := oss.SizeUnit(this.total)
		return fmt.Sprintf("%0.1f%s/%0.1f%s", currentNum, currentUnit, totalNum, totalUnit)
	})
}

func (this *base) speed(key string, size int64, expiration time.Duration, f func(float64) string) string {

	timeKey := "time_" + key
	cacheKey := "speed_" + key
	//最后的数据时间
	lastTime, _ := this.cache.GetOrSetByHandler(timeKey, func() (interface{}, error) {
		return time.Time{}, nil
	})

	//记录这次时间,用于下次计算时间差
	now := time.Now()
	this.cache.Set(timeKey, now)

	//尝试从缓存获取速度,存在则直接返回,由expiration控制
	if val, ok := this.cache.Get(cacheKey); ok {
		return val.(string)
	}

	//计算速度
	size = conv.SelectInt64(size >= 0, size, 0)
	spendSize := float64(size) / now.Sub(lastTime.(time.Time)).Seconds()
	s := f(spendSize)
	this.cache.Set(cacheKey, s, expiration)
	return s
}

func (this *base) Speed() Element {
	return ElementFunc(func() string {
		return this.speed("Speed", this.last, time.Millisecond*500, func(size float64) string {
			return fmt.Sprintf("%0.1f/s", size)
		})
	})
}

func (this *base) SpeedUnit() Element {
	return ElementFunc(func() string {
		return this.speed("SpeedUnit", this.last, time.Millisecond*500, func(size float64) string {
			f, unit := oss.SizeUnit(int64(size))
			return fmt.Sprintf("%0.1f%s/s", f, unit)
		})
	})
}

func (this *base) SpeedAvg() Element {
	return ElementFunc(func() string {
		speedSize := float64(this.current) / time.Since(this.startTime).Seconds()
		return fmt.Sprintf("%0.1f/s", speedSize)
	})
}

func (this *base) SpeedUnitAvg() Element {
	return ElementFunc(func() string {
		speedSize := float64(this.current) / time.Since(this.startTime).Seconds()
		f, unit := oss.SizeUnit(int64(speedSize))
		return fmt.Sprintf("%0.1f%s/s", f, unit)
	})
}

func (this *base) Used() Element {
	return ElementFunc(func() string {
		return time.Now().Sub(this.startTime).String()
	})
}

func (this *base) UsedSecond() Element {
	return ElementFunc(func() string {
		return fmt.Sprintf("%0.1fs", time.Now().Sub(this.startTime).Seconds())
	})
}

func (this *base) Remain() Element {
	return ElementFunc(func() string {
		rate := float64(this.current) / float64(this.total)
		spend := time.Now().Sub(this.startTime)
		remain := "0s"
		if rate > 0 {
			sub := time.Duration(float64(spend)/rate - float64(spend))
			remain = (sub - sub%time.Second).String()
		}
		return remain
	})
}

/*



 */

var (
	// DefaultClient 默认客户端,下载大文件的时候需要设置长的超时时间
	DefaultClient = http.NewClient().SetTimeout(0)
)

func (this *base) DownloadHTTP(source, filename string, proxy ...string) (int64, error) {
	if err := DefaultClient.SetProxy(conv.GetDefaultString("", proxy...)); err != nil {
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
	total := int64(0)
	buf := make([]byte, bufSize)
	for {
		n, err := buff.Read(buf)
		if err != nil && err != io.EOF {
			return total, err
		}
		total += int64(n)
		this.Add(int64(n))
		this.Flush()
		if _, err := w.Write(buf[:n]); err != nil {
			return total, err
		}
		if err == io.EOF {
			return total, nil
		}
	}
}
