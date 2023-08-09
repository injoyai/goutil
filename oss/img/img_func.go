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

// DrawImg 增加水印
func DrawImg(src, mark image.Image, offset image.Point) (image.Image, error) {
	srcBounds := src.Bounds()
	newImg := image.NewNRGBA(srcBounds)
	draw.Draw(newImg, srcBounds, src, image.ZP, draw.Src)
	draw.Draw(newImg, mark.Bounds().Add(offset), mark, image.ZP, draw.Over)
	return newImg, nil
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
