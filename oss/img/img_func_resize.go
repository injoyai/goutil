package img

import (
	"bytes"
	"image"
	"io"

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

func ResizeReader(r io.Reader, maxSize uint) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return Resize(img, maxSize), nil
}

func ResizeBytes(bs []byte, maxSize uint) (image.Image, error) {
	return ResizeReader(bytes.NewReader(bs), maxSize)
}
