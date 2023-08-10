package main

import (
	"encoding/base64"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/oss/img"
	"github.com/injoyai/goutil/oss/img/imgcode"
	"github.com/injoyai/goutil/oss/img/qrcode"
	"github.com/injoyai/goutil/str"
	"github.com/injoyai/logs"
)

func main() {

	{
		im := img.New(100, 100, img.Green)
		err := im.Save("./oss/img/testdata/img.png")
		logs.PrintErr(err)
	}
	{
		err := qrcode.EncodeLocal("./oss/img/testdata/qr.png", "666")
		logs.PrintErr(err)
	}
	{
		_, baseImg, err := imgcode.Get()
		logs.PrintErr(err)
		baseImg = str.CropFirst(baseImg, ",", false)
		bs, err := base64.StdEncoding.DecodeString(baseImg)
		logs.PrintErr(err)
		oss.New("./oss/img/testdata/imgcode1.png", bs)
	}
	{
		_, baseImg, err := imgcode.Get()
		logs.PrintErr(err)
		baseImg = str.CropFirst(baseImg, ",", false)
		bs, err := base64.StdEncoding.DecodeString(baseImg)
		logs.PrintErr(err)
		oss.New("./oss/img/testdata/imgcode2.png", bs)
	}
}
