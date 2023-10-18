/*****************************************************************************
*名称:	时间工具包
*功能:	获取时间的节点以及时间的时间段加减操作
*作者:	钱纯净
******************************************************************************/

package times

import (
	"time"
)

type (
	Duration = time.Duration
	Month    = time.Month
)

const (
	FormatDefault = "2006-01-02 15:04:05"
	FormatDate    = "2006-01-02"
	FormatTime    = "15:04:05"

	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
	Day         = Hour * 24
	Week        = Day * 7
)

type Time struct {
	time.Time
}

// Now 取当前时间
func Now() Time {
	return Time{
		Time: time.Now(),
	}
}

// New 新建时间,秒
func New(sec int64) Time {
	return NewSec(sec)
}

func NewSec(sec int64) Time {
	return NewNano(sec * 1e9)
}

func NewMill(mill int64) Time {
	return NewNano(mill * 1e6)
}

func NewMicro(micro int64) Time {
	return NewNano(micro * 1e3)
}

func NewNano(nano int64) Time {
	return Time{
		Time: time.Unix(nano/1e9, nano%1e9),
	}
}

// Parse string类型转time
func Parse(format, value string) (Time, error) {
	t, err := time.Parse(format, value)
	return Time{Time: t}, err
}

// ParseDefault string类型转time
func ParseDefault(s string) (Time, error) {
	t, err := time.Parse(FormatDefault, s)
	return Time{Time: t}, err
}

func Date(year, month, day, hour, min, second int) Time {
	return Time{Time: time.Date(year, time.Month(month), day, hour, min, second, 0, time.Local)}
}

func (this Time) String() string {
	return this.Format(FormatDefault)
}

// Format 转为字符串,格式比如 "2006-01-02 15:04:05"
func (this Time) Format(format string) string {
	return this.Time.Format(format)
}

// Date 年月日秒
func (this Time) Date() (int, int, int, int) {
	year, month, day := this.Time.Date()
	return year, int(month), day, this.Second()
}

// Year 年
func (this Time) Year() int {
	return this.Time.Year()
}

// Quarter 1-4季,当年
func (this Time) Quarter() int {
	return (this.Month()-1)/3 + 1
}

// Month 1-12月,当年
func (this Time) Month() int {
	return int(this.Time.Month())
}

// Day 1-31日,当月
func (this Time) Day() int {
	return this.Time.Day()
}

// Weekday 0-6星期,星期天是0
func (this Time) Weekday() int {
	return int(this.Time.Weekday())
}

// Hour 0-23时,当天
func (this Time) Hour() int {
	return this.Time.Hour()
}

// Minute 0-59分,时
func (this Time) Minute() int {
	return this.Time.Minute()
}

// Second 0-59秒,分
func (this Time) Second() int {
	return this.Time.Second()
}

func (this Time) Duration() time.Duration {
	return time.Duration(this.UnixNano())
}

//********************************************分界线**************************************************//

func (this Time) Add(d time.Duration) Time {
	this.Time = this.Time.Add(d)
	return this
}

// AddSec 秒,加减
func (this Time) AddSec(sec int) Time {
	return this.Add(time.Second * time.Duration(sec))
}

// AddMin 分,加减
func (this Time) AddMin(minute int) Time {
	return this.Add(time.Minute * time.Duration(minute))
}

// AddHour 小时加减
func (this Time) AddHour(hour int) Time {
	return this.Add(time.Hour * time.Duration(hour))
}

// AddDay 天,加减
func (this Time) AddDay(day int) Time {
	this.Time = this.AddDate(0, 0, day)
	return this
}

// AddWeek 周,加减
func (this Time) AddWeek(week int) Time {
	this.Time = this.AddDate(0, 0, week*7)
	return this
}

// AddMonth 月,加减
func (this Time) AddMonth(month int) Time {
	this.Time = this.AddDate(0, month, 0)
	return this
}

// AddYear 年,加减
func (this Time) AddYear(year int) Time {
	this.Time = this.AddDate(year, 0, 0)
	return this
}

//************************************************分界线**********************************************//

// IntegerSecond 取整秒
func (this Time) IntegerSecond() Time {
	return this.Add(-this.Duration() % Second)
}

// IntegerMinute 取整分
func (this Time) IntegerMinute() Time {
	return this.Add(-this.Duration() % Minute)
}

// IntegerHour 取整点
func (this Time) IntegerHour() Time {
	return this.Add(-this.AddHour(8).Duration() % Hour)
}

// IntegerDay 取整天
func (this Time) IntegerDay() Time {
	return this.Add(-this.AddHour(8).Duration() % (Day))
}

// IntegerWeek 取整周,周一
func (this Time) IntegerWeek() Time {
	return this.AddDay(-this.AddHour(8).Weekday() + 1).IntegerDay()
}

// IntegerMonth 取整月
func (this Time) IntegerMonth() Time {
	year, month, _, _ := this.Date()
	this.Time = time.Date(year, Month(month), 1, 0, 0, 0, 0, time.Local)
	return this
}

func (this Time) IntegerQuarter() Time {
	year := this.Year()
	this.Time = time.Date(year, Month((this.Quarter()-1)*3), 1, 0, 0, 0, 0, time.Local)
	return this
}

// IntegerYear 取整年
func (this Time) IntegerYear() Time {
	year := this.Year()
	this.Time = time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	return this
}

//********************************************************分界线**********************************************************//

// Sub 差值
func (this Time) Sub() time.Duration {
	return time.Now().Sub(this.Time)
}
