package img

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

// Resize
// 传0等比例压缩 1153k >>> 288k
// 压缩和缩放,指定最大变,等比例
func Resize(img image.Image, maxSize uint) image.Image {
	oriBounds := img.Bounds()
	oriWidth := uint(oriBounds.Dx())
	oriHeight := uint(oriBounds.Dy())
	if maxSize == 0 {
		maxSize = oriWidth
		if oriHeight > maxSize {
			maxSize = oriHeight
		}
	}
	return resize.Thumbnail(maxSize, maxSize, img, resize.Lanczos3)
}

func ResizeBytesJpeg(bs []byte, maxSize uint) ([]byte, error) {
	i, err := jpeg.Decode(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	w := bytes.NewBuffer(nil)
	err = jpeg.Encode(w, Resize(i, maxSize), nil)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func ResizeBytesPng(bs []byte, maxSize uint) ([]byte, error) {
	i, err := png.Decode(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	w := bytes.NewBuffer(nil)
	err = png.Encode(w, Resize(i, maxSize))
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
