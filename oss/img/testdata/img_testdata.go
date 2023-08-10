package main

import (
	"github.com/injoyai/goutil/oss/img"
	"github.com/injoyai/logs"
)

func main() {
	im := img.New(100, 100, img.Green)
	err := im.Save("./oss/img/testdata/test.png")
	logs.PrintErr(err)
}
