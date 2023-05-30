package idcard

import (
	"github.com/injoyai/conv"
	"time"
)

type Entity struct {
	Province string //省
	City     string //市
	Area     string //区
	Year     int    //年
	Month    int    //月
	Day      int    //日
	Sex      string //性别
	IDCard   string //身份证
}

// Age 获取年龄周岁
func (this *Entity) Age() int {
	if this.Year == 0 {
		return -1
	}
	year, month, day := time.Now().Date()
	age := year - this.Year - 1
	if int(month)-this.Month > 0 || (int(month)-this.Month == 0 && day-this.Day >= 0) {
		age++
	}
	return age
}

// New 身份证数据
func New(idcard string) *Entity {
	data := &Entity{
		IDCard: idcard,
		Sex:    "未知",
	}
	if len(idcard) <= 18 {

		if len(idcard) >= 2 {
			data.Province = SiteMap[idcard[:2]+"0000"]
		}
		if len(idcard) >= 4 {
			data.City = SiteMap[idcard[:4]+"00"]
		}
		if len(idcard) >= 6 {
			data.Area = SiteMap[idcard[:6]]
		}
		if len(idcard) >= 10 {
			data.Year = conv.Int(idcard[6:10])
		}
		if len(idcard) >= 12 {
			data.Month = conv.Int(idcard[10:12])
		}
		if len(idcard) >= 14 {
			data.Day = conv.Int(idcard[12:14])
		}
		if len(idcard) >= 17 {
			data.Sex = "男"
			if conv.Int(idcard[16:17])%2 == 0 {
				data.Sex = "女"
			}
		}
	}
	return data
}
