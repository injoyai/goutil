package main

import (
	"github.com/fatih/color"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/str/bar"
	"github.com/injoyai/logs"
	"os"
)

func init() {
	logs.DefaultErr.SetWriter(logs.Stdout)
}

func main() {

	bar.Demo()

	url := "https://github.com/injoyai/downloader/releases/latest/download/downloader.exe"
	filename := "./downloader.exe"
	b := bar.New(0)
	b.SetColor(color.BgCyan)
	b.SetStyle('#')
	for {
		err := b.DownloadHTTP(url, filename)
		if !logs.PrintErr(err) {
			os.Remove(filename)
			break
		}
	}
	g.Input("请按回车键退出...")
}
