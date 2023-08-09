package img

import (
	"bytes"
	"github.com/injoyai/goutil/oss"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// New 新建
func New(wide, high int, backColor ...color.Color) image.Image {
	alpha := image.NewNRGBA(image.Rect(0, 0, wide, high))
	if len(backColor) > 0 {
		for x := 0; x < wide; x++ {
			for y := 0; y < high; y++ {
				alpha.Set(x, y, backColor[0])
			}
		}
	}
	return alpha
}

// DrawImg 画上图片
func DrawImg(img1, img2 image.Image, offset ...image.Point) (image.Image, error) {
	srcBounds := img1.Bounds()
	newImg := image.NewNRGBA(srcBounds)
	draw.Draw(newImg, srcBounds, img1, image.ZP, draw.Src)
	rectangle := img2.Bounds()
	if len(offset) > 0 {
		rectangle.Add(offset[0])
	}
	draw.Draw(newImg, rectangle, img2, image.ZP, draw.Over)
	return newImg, nil
}

// JoinImg 拼接图片
func JoinImg(img1, img2 image.Image, offsets ...image.Point) (img image.Image, err error) {
	offset := image.Point{}
	if len(offsets) > 0 {
		offset = offsets[0]
	}
	wide, high := img1.Bounds().Dx(), img1.Bounds().Dy()
	if img2.Bounds().Max.X+offset.X > wide {
		wide = img2.Bounds().Max.X + offset.X
	}
	if img2.Bounds().Max.Y+offset.Y > high {
		high = img2.Bounds().Max.Y + offset.Y
	}
	img = New(wide, high)
	img, err = DrawImg(img, img1)
	if err != nil {
		return nil, err
	}
	img, err = DrawImg(img, img2, offset)
	return
}

func Open(filename string) (image.Image, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	return image.Decode(file)
}

func OpenJpeg(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(file)
}

func OpenPng(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return png.Decode(file)
}

func Save(filename string, img image.Image) error {
	switch strings.ToUpper(filepath.Ext(filename)) {
	case "PNG":
		return SavePng(filename, img)
	case "JPG", "JPEG":
		return SaveJpeg(filename, img)
	}
	return SavePng(filename, img)
}

func SaveJpeg(path string, img image.Image) error {
	bs := bytes.NewBuffer(nil)
	if err := jpeg.Encode(bs, img, nil); err != nil {
		return err
	}
	return oss.New(path, bs)
}

func SavePng(path string, img image.Image) error {
	bs := bytes.NewBuffer(nil)
	if err := png.Encode(bs, img); err != nil {
		return err
	}
	return oss.New(path, bs)
}
