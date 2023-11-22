package oss

import (
	"fmt"
	"github.com/injoyai/conv"
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

// SizeString 字节数量字符表现方式
func SizeString(size uint64, decimal ...int) string {
	base := 0
	for base = 0; base <= 130; base += 10 {
		if size < 1<<(base+10) {
			break
		}
	}
	unit := mapSizeUnit[base]
	f := float64(size) / float64(int(1)<<base)
	d := conv.GetDefaultInt(2, decimal...)
	return fmt.Sprintf(fmt.Sprintf("%%.%df%%s", d), f, unit)
}
