package img

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var defaultBrush *Brush

// Brush 画笔
type Brush struct {
	Type  *truetype.Font //字体样式
	Size  int            //字体大小
	Color color.Color    //字体颜色
}

// SetSize 设置画笔大小
func (this *Brush) SetSize(size int) {
	this.Size = size
}

// SetColor 设置画笔颜色
func (this *Brush) SetColor(color color.Color) {
	this.Color = color
}

// SetTypeLocal 设置画笔样式本地
func (this *Brush) SetTypeLocal(filename string) error {
	font, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return this.SetType(font)
}

// SetType 设置画笔样式
func (this *Brush) SetType(font []byte) error {
	fontType, err := truetype.Parse(font)
	if err != nil {
		return err
	}
	this.Type = fontType
	return nil
}

// DrawText 写入文字
func (this *Brush) DrawText(img draw.Image, content string, offsets ...image.Point) error {
	c := freetype.NewContext()
	c.SetDPI(45)
	c.SetFont(this.Type)
	c.SetHinting(font.HintingFull)
	c.SetFontSize(float64(this.Size))
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(this.Color))
	offset := image.Point{}
	if len(offsets) > 0 {
		offset = offsets[0]
	}
	_, err := c.DrawString(content, freetype.Pt(offset.X, offset.Y))
	return err
}

//========================= 图片绘画 =========================//

// DrawPoint 画点
func (this *Brush) DrawPoint(img draw.Image, site image.Point) {
	for x := site.X - this.Size; x < site.X+this.Size; x++ {
		for y := site.Y - this.Size; y < site.Y+this.Size; y++ {
			distance := math.Sqrt(math.Pow(float64(x-site.X), 2) + math.Pow(float64(y-site.Y), 2))
			if distance < float64(this.Size)-1 {
				img.Set(x, y, this.Color)
			} else if distance < float64(this.Size) {
				//TODO 锯齿处理
				img.Set(x, y, this.Color)
			}
		}
	}
}

// DrawLine 画线
func (this *Brush) DrawLine(img draw.Image) {

}

// DrawCircle 画圆
func (this *Brush) DrawCircle(img draw.Image, centre image.Point, radii int) {
	for x := centre.X - radii; x < centre.X+radii; x++ {
		for y := centre.Y - radii; y < centre.Y+radii; y++ {
			distance := math.Sqrt(math.Pow(float64(x-centre.X), 2) + math.Pow(float64(y-centre.Y), 2))
			if distance >= float64(radii)-1 && distance < float64(radii) {
				img.Set(x, y, this.Color)
			}
		}
	}
}

// DrawRectangle 画正方形
func (this *Brush) DrawRectangle(img draw.Image, start image.Point, wide, high int) {
	for x := start.X; x <= start.X+wide; x++ {
		img.Set(x, start.Y, this.Color)
	}
	for x := start.X; x <= start.X+wide; x++ {
		img.Set(x, start.Y+high, this.Color)
	}
	for y := start.Y; y <= start.Y+high; y++ {
		img.Set(start.X, y, this.Color)
	}
	for y := start.Y; y <= start.Y+high; y++ {
		img.Set(start.X+wide, y, this.Color)
	}
}

//========================= 新建画笔 =========================//

// NewBrush 新建画笔
func NewBrush() (*Brush, error) {
	//加载字体
	fontType, err := truetype.Parse(defFontBytes)
	if err != nil {
		return nil, err
	}
	return &Brush{
		Type:  fontType,
		Size:  50,
		Color: color.Black,
	}, nil
}
