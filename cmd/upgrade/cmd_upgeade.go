package main

import (
	"github.com/injoyai/base/oss"
	"github.com/injoyai/goutil/cmd/nac"
	"github.com/injoyai/goutil/string/bar"
	"github.com/injoyai/logs"
)

func main() {
	nac.Init()
	url := "https://github.com/injoyai/goutil/raw/main/cmd/in.exe"
	logs.PrintErr(bar.Download(url, "./in.exe"))
	oss.Input("升级成功,按任意键退出...")
}
