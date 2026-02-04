package oss

import (
	"fmt"
	"math"
	"strings"

	"github.com/injoyai/conv"
)

const (
	B  = 1
	KB = B * 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
	PB = TB * 1024
	EB = PB * 1024
)

type Volume uint64

func (this Volume) Uint64() uint64 {
	return uint64(this)
}

func (this Volume) String() string {
	size, unit := this.SizeUnit()
	return fmt.Sprintf("%.2f%s", size, unit)
}

func (this Volume) SizeUnit() (float64, string) {
	if this == 0 {
		return 0, mapSizeUnit[0]
	}

	f := float64(this)
	i := 0

	for f >= 1024 && i < len(mapSizeUnit)-1 {
		f /= 1024
		i++
	}

	return f, mapSizeUnit[i]
}

// ParseVolume 解析体积,
func ParseVolume(s string) Volume {
	total := Volume(0)
	size := ""
	unit := ""
	hasUnit := false

	add := func() {
		if hasUnit {
			for i, u := range mapSizeUnit {
				if strings.ToUpper(unit) == u {
					total += Volume(conv.Float64(size) * math.Pow(1024, float64(i)))
				}
			}
			hasUnit = false
			size = ""
			unit = ""
		}
	}

	for _, v := range s {
		if (v >= '0' && v <= '9') || v == '.' {
			add()
			size += string(v)
		} else {
			unit += string(v)
			hasUnit = true
		}
	}
	add()

	return total
}
