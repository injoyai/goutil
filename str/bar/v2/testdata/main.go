package main

import (
	"fmt"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/str/bar/v2"
	"github.com/injoyai/ios"
	"github.com/injoyai/logs"
	"os"
	"time"
)

func main() {
	bar.Demo()

	{
		logs.Debug("文件示例:")
		url := "https://github.com/injoyai/downloader/releases/latest/download/downloader.exe"
		filename := "./downloader.exe"
		b := bar.New(bar.WithDefaultFormat(func(p bar.Plan) {
			p.SetColor(bar.BgYellow)
			p.SetStyle('#')
			p.SetPrefix("(")
			p.SetSuffix(")")
		}))

		b.SetWriter(ios.WriteFunc(func(p []byte) (int, error) {
			return fmt.Print(string(p))
		}))
		for {
			_, err := b.DownloadHTTP(url, filename, "http://127.0.0.1:1081")
			if !logs.PrintErr(err) {
				os.Remove(filename)
				break
			}
			<-time.After(time.Second * 5)
		}
	}

	{
		logs.Debug("失败示例:")
		b := bar.New(bar.WithTotal(100))
		go func(b bar.Bar) {
			<-time.After(time.Second * 10)
			b.Close()
		}(b)
		b.Add(1)
		b.Flush()
		<-b.Done()
	}

	g.Input("请按回车键退出...")
}
