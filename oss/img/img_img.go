package img

import (
	"bytes"
	"image"
	"image/png"
)

type Img struct {
	image.Image
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

func (this *Img) Resize(maxSize uint) {
	this.Image = Resize(this.Image, maxSize)
}

func (this *Img) DrawImg(img image.Image, offset ...image.Point) (err error) {
	this.Image, err = DrawImg(this.Image, img, offset...)
	return err
}

func (this *Img) JoinImg(img image.Image, offset ...image.Point) (err error) {
	this.Image, err = JoinImg(this.Image, img, offset...)
	return err
}

//func (this *Img) DrawText(text string) {
//	this.Image = Resize(maxSize, this.Image)
//}

// Buffer 返回字节buffer
func (this *Img) Buffer() (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, this.Image)
	return buf, err
}
