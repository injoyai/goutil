package times

import (
	"testing"
	"time"
)

func TestQuarter(t *testing.T) {
	for i := 1; i < 13; i++ {
		n := Quarter(time.Date(2020, Month(i), 1, 0, 0, 0, 0, time.Local))
		switch i {
		case 1, 2, 3:
			if n != 1 {
				t.Error("季度错误")
				return
			}
		case 4, 5, 6:
			if n != 2 {
				t.Error("季度错误")
				return
			}
		case 7, 8, 9:
			if n != 3 {
				t.Error("季度错误")
				return
			}
		case 10, 11, 12:
			if n != 4 {
				t.Error("季度错误")
				return
			}
		}
	}
}

func TestInteger(t *testing.T) {
	now := time.Now() //.AddDate(0, 4, 4)
	t.Log("秒: ", IntegerSecond(now).UnixNano())
	t.Log("秒: ", IntegerSecond(now).String())
	t.Log("分: ", IntegerMinute(now).String())
	t.Log("时: ", IntegerHour(now).String())
	t.Log("天: ", IntegerDay(now).String())
	t.Log("周: ", IntegerWeek(now).String())
	t.Log("月: ", IntegerMonth(now).String())
	t.Log("季: ", IntegerQuarter(now).String())
	t.Log("年: ", IntegerYear(now).String())
}
