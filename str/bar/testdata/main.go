package main

import (
	"github.com/fatih/color"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/str/bar"
	"github.com/injoyai/logs"
	"os"
	"time"
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
		err := b.DownloadHTTP(url, filename) //, "http://127.0.0.1:1081")
		if !logs.PrintErr(err) {
			os.Remove(filename)
			break
		}
		<-time.After(time.Second * 5)
	}

	{
		logs.Debug("失败示例:")
		b = bar.New(100)
		go func(b bar.Interface) {
			<-time.After(time.Second * 10)
			b.Close()
		}(b)
		b.Run()
	}

	g.Input("请按回车键退出...")
}
