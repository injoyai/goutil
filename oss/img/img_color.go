package img

import "image/color"

type Color struct {
	color.Color
}

// Avg 颜色平均
func (this Color) Avg(color color.Color) Color {
	r, g, b, a := this.Color.RGBA()
	r2, g2, b2, a2 := color.RGBA()
	return Color{NewRGBA(uint8((r+r2)/2), uint8((g+g2)/2), uint8((b+b2)/2), uint8((a+a2)/2))}
}

func NewRGB(r, g, b uint8) color.Color {
	return color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}

func NewRGBA(r, g, b, a uint8) color.Color {
	return color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}
