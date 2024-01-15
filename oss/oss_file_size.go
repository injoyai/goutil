package oss

import (
	"fmt"
	"github.com/injoyai/conv"
	"time"
)

var mapSizeUnit = map[int]string{
	0:   "B",
	10:  "KB",
	20:  "MB",
	30:  "GB",
	40:  "TB",
	50:  "PB",
	60:  "EB",
	70:  "ZB",
	80:  "YB",
	90:  "BB",
	100: "NB",
	110: "DB",
	120: "CB",
	130: "XB",
}

func Size(b int64) (float64, string) {
	for n := 0; n <= 130; n += 10 {
		if b < 1<<(n+10) {
			if n == 0 {
				return float64(b), mapSizeUnit[n]
			}
			return float64(b) / float64(int64(1)<<n), mapSizeUnit[n]
		}
	}
	return float64(b), mapSizeUnit[0]
}

// SizeString 字节数量字符表现方式
func SizeString(b int64, decimal ...int) string {
	size, unit := Size(b)
	d := conv.GetDefaultInt(1, decimal...)
	return fmt.Sprintf(fmt.Sprintf("%%.%df%%s", d), size, unit)
}

func SizeSpend(b int64, d time.Duration) (float64, string) {
	size, unit := Size(b)
	return size / d.Seconds(), unit
}

func SizeSpendString(b int64, sub time.Duration, decimal ...int) string {
	spend, unit := SizeSpend(b, sub)
	d := conv.GetDefaultInt(1, decimal...)
	return fmt.Sprintf(fmt.Sprintf("%%.%df%%s/s", d), spend, unit)
}
