/*****************************************************************************
*名称:	时间工具包
*功能:	获取时间的节点以及时间的时间段加减操作
*作者:	钱纯净
******************************************************************************/

package times

import (
	"fmt"
	"github.com/injoyai/conv"
	"strings"
	"time"
)

type times struct {
	time int64 //秒
	//nanosecond int64 //纳秒,纳秒太快,加计算控制不住时间的准确性,基础误差0.03-0.01秒
	start time.Time //用于打印输出用时
}

// Now 取当前时间
func Now() times {
	t := time.Now()
	return times{
		time: t.Unix(),
		//Nanosecond: t.UnixNano(),
		start: t,
	}
}

// New 新建时间,秒
func New(n int64) times {
	return times{
		time:  n,
		start: time.Unix(n, 0),
	}
}

// Unix 对应time.Unix()
func Unix(n int64) times {
	return times{
		time:  n,
		start: time.Unix(n, 0),
	}
}

// String string类型转time
func String(s string) (times, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		return times{}, err
	}
	return times{
		time:  t.Unix(),
		start: t,
	}, nil
}

func Date(year, month, day, hour, min, second int) times {
	return times{time: time.Date(year, time.Month(month), day, hour, min, second, 0, time.Local).Unix()}
}

// Unix 转为时间戳
func (this times) Unix() int64 {
	return this.time
}

// String 转为字符串,格式比如 "年-月-日 时-分-秒" "2006-01-02 15:04:05"
func (this times) String(s ...string) string {
	var str string

	if len(s) > 0 {
		str = strings.Replace(s[0], "年", "2006", -1)
		str = strings.Replace(str, "月", "01", -1)
		str = strings.Replace(str, "日", "02", -1)
		str = strings.Replace(str, "时", "15", -1)
		str = strings.Replace(str, "分", "04", -1)
		str = strings.Replace(str, "秒", "05", -1)
		str = time.Unix(this.time, 0).Format(str)
	} else {
		//0.4秒
		str = time.Unix(this.time, 0).Format("2006-01-02 15:04:05")
	}
	return str
}

// Date 年月日秒
func (this times) Date() (int, int, int, int) {
	year, month, day := this.ToTime().Date()
	second := this.Unix() - this.IntegerDay().Unix()
	return year, int(month), day, int(second)
}

// Year 年
func (this times) Year() int {
	return this.ToTime().Year()
}

// Quarter 1-4季,当年
func (this times) Quarter() int {
	return (this.Month()-1)/3 + 1
}

// Month 1-12月,当年
func (this times) Month() int {
	return int(this.ToTime().Month())
}

// Day 1-31日,当月
func (this times) Day() int {
	return this.ToTime().Day()
}

// Week 0-6星期,星期天是0
func (this times) Week() int {
	return int(this.ToTime().Weekday())
}

// Hour 0-23时,当天
func (this times) Hour() int {
	return int((this.time - New(this.time).IntegerDay().Unix()) / (60 * 60))
}

// Minute 0-59分,时
func (this times) Minute() int {
	return int(this.time-New(this.time).IntegerHour().Unix()) / 60
}

// Second 0-59秒,分
func (this times) Second() int64 {
	return this.Unix() - this.IntegerMin().Unix()
}

//********************************************分界线**************************************************//

// ToTime 变成标准库的time
func (this times) ToTime() time.Time {
	return time.Unix(this.time, 0)
}

// Add 秒,加减
func (this times) Add(a int) times {
	this.time += int64(a)
	return this
}

// AddSec 秒,加减
func (this times) AddSec(a int) times {
	this.time += int64(a)
	return this
}

// AddMin 分,加减
func (this times) AddMin(a int) times {
	this.time += int64(a) * 60
	return this
}

// AddHour 小时加减
func (this times) AddHour(a int) times {
	this.time += int64(a) * 60 * 60
	return this
}

// AddDay 天,加减
func (this times) AddDay(a int) times {
	this.time += int64(a) * 60 * 60 * 24
	return this
}

// AddWeek 周,加减
func (this times) AddWeek(a int) times {
	this.time = time.Unix(this.time, 0).AddDate(0, 0, a*7).Unix()
	return this
}

// AddMonth 月,加减
func (this times) AddMonth(a int) times {
	this.time = time.Unix(this.time, 0).AddDate(0, a, 0).Unix()
	return this
}

// AddYear 年,加减
func (this times) AddYear(a int) times {
	this.time = time.Unix(this.time, 0).AddDate(a, 0, 0).Unix()
	return this
}

//************************************************分界线**********************************************//

// IntegerMin 取整分
func (this times) IntegerMin() times {
	this.time = this.time - this.time%60
	return this
}

// IntegerHour 取整点
func (this times) IntegerHour() times {
	this.time = this.time - this.time%(60*60)
	return this
}

// IntegerDay 取整天
func (this times) IntegerDay() times {
	this.time = this.time - (this.time+60*60*8)%(60*60*24)
	return this
}

// IntegerWeek 取整周,周一
func (this times) IntegerWeek() times {
	t := int(time.Unix(this.time, 0).Weekday())
	this.time = this.AddDay(-t + 1).IntegerDay().Unix()
	return this
}

// IntegerMonth 取整月
func (this times) IntegerMonth() times {
	year, month, _ := time.Unix(this.time, 0).Date()
	this.time = time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Unix()
	return this
}

// IntegerYear 取整年
func (this times) IntegerYear() times {
	year := time.Unix(this.time, 0).Year()
	this.time = time.Date(year, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	return this
}

//********************************************************分界线**********************************************************//

// SubSecond 差值秒
func (this times) SubSecond() int64 {
	return int64(this.Sub() / time.Second)
}

// Sub 差值
func (this times) Sub() time.Duration {
	return time.Now().Sub(this.start)
}

func (this times) PrintSub(s ...string) {
	x := conv.GetDefaultString("耗时", s...)
	fmt.Printf("[%s] %v\n", x, this.Sub())
}
