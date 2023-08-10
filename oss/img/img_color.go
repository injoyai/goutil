package img

import "image/color"

var (
	White   = NewRGB(255, 255, 255) //白色
	Black   = NewRGB(0, 0, 0)       //黑色
	Red     = NewRGB(255, 0, 0)     //红色
	Green   = NewRGB(0, 255, 0)     //绿色
	Blue    = NewRGB(0, 0, 255)     //蓝色
	Cyan    = NewRGB(0, 255, 255)   //青色
	Magenta = NewRGB(255, 0, 255)   //品红
	Grey    = NewRGB(192, 192, 192) //灰色
)

func NewRGB(r, g, b uint8) color.Color {
	return NewRGBA(r, g, b, 255)
}

func NewRGBA(r, g, b, a uint8) color.Color {
	return color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func ColorAvg(colors ...color.Color) color.Color {
	var r, g, b, a, length uint64
	for _, c := range colors {
		r0, g0, b0, a0 := c.RGBA()
		r += uint64(r0)
		g += uint64(g0)
		b += uint64(b0)
		a += uint64(a0)
	}
	return NewRGBA(uint8(r/length), uint8(g/length), uint8(b/length), uint8(a/length))
}
