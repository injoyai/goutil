package img

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"io"
)

type (
	Image = image.Image
	Color = color.Color
)

func New(wide, high int, backColor ...color.Color) *Img {
	alpha := image.NewNRGBA(image.Rect(0, 0, wide, high))
	if len(backColor) > 0 {
		for x := 0; x < wide; x++ {
			for y := 0; y < high; y++ {
				alpha.Set(x, y, backColor[0])
			}
		}
	}
	return &Img{alpha}
}

type Img struct {
	draw.Image
}

func (this *Img) Resize(maxSize uint) image.Image {
	return Resize(this.Image, maxSize)
}

func (this *Img) DrawImg(img image.Image, offset ...image.Point) (err error) {
	this.Image, err = DrawImg(this.Image, img, offset...)
	return err
}

func (this *Img) DrawReader(r io.Reader, offset ...image.Point) (err error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	this.Image, err = DrawImg(this.Image, img, offset...)
	return err
}

func (this *Img) DrawBytes(bs []byte, offset ...image.Point) (err error) {
	return this.DrawReader(bytes.NewReader(bs), offset...)
}

func (this *Img) JoinImg(img image.Image, offset ...image.Point) (err error) {
	this.Image, err = JoinImg(this.Image, img, offset...)
	return err
}

func (this *Img) JoinReader(r io.Reader, offset ...image.Point) (err error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	this.Image, err = JoinImg(this.Image, img, offset...)
	return err
}

func (this *Img) JoinBytes(bs []byte, offset ...image.Point) (err error) {
	return this.JoinReader(bytes.NewReader(bs), offset...)
}

func (this *Img) DrawText(text string, offset ...image.Point) (err error) {
	if defaultBrush == nil {
		defaultBrush, err = NewBrush()
		if err != nil {
			return err
		}
	}
	return defaultBrush.DrawText(this.Image, text, offset...)
}

// BufferPng 返回字节buffer
func (this *Img) BufferPng() (*bytes.Buffer, error) {
	return BufferPng(this.Image)
}

// BufferJpeg 返回字节buffer
func (this *Img) BufferJpeg() (*bytes.Buffer, error) {
	return BufferJpeg(this.Image)
}

func (this *Img) Save(filename string) error {
	return Save(filename, this.Image)
}

// SavePng 保存成png
func (this *Img) SavePng(filename string) error {
	return SavePng(filename, this.Image)
}

// SaveJpeg 保存成jpeg
func (this *Img) SaveJpeg(filename string) error {
	return SaveJpeg(filename, this.Image)
}
