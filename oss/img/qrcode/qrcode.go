package qrcode

import (
	"bytes"
	"github.com/injoyai/goutil/oss"
	"github.com/skip2/go-qrcode"
	qrcode2 "github.com/tuotoo/qrcode"
	"image"
	"image/color"
	"image/png"
	"io"
)

func Encode(text string, color ...color.Color) (image.Image, error) {
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	if len(color) > 0 {
		qr.ForegroundColor = color[0]
	}
	return qr.Image(256), nil
}

func EncodeBuffer(text string, color ...color.Color) (*bytes.Buffer, error) {
	img, err := Encode(text, color...)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	err = png.Encode(buf, img)
	return buf, err
}

func EncodeBytes(text string, color ...color.Color) ([]byte, error) {
	buf, err := EncodeBuffer(text, color...)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodeLocal(filename, text string, color ...color.Color) error {
	bs, err := EncodeBytes(text, color...)
	if err != nil {
		return err
	}
	return oss.New(filename, bs)
}

func Decode(img image.Image) (string, error) {
	matrix, err := qrcode2.DecodeImg(img, "")
	if err != nil {
		return "", err
	}
	return matrix.Content, nil
}

func DecodeReader(r io.Reader) (string, error) {
	matrix, err := qrcode2.Decode(r)
	if err != nil {
		return "", err
	}
	return matrix.Content, nil
}

func DecodeBytes(bs []byte) (string, error) {
	return DecodeReader(bytes.NewReader(bs))
}
