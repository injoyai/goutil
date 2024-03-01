package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/str/bar"
	"github.com/injoyai/io"
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
	b.SetColor(color.BgYellow)
	b.SetStyle('#')
	b.SetWriter(io.WriteFunc(func(p []byte) (int, error) {
		return fmt.Print(string(p))
	}))
	b.AddOption(func(f *bar.Format) {
		f.Bar.SetPrefix("(")
		f.Bar.SetSuffix(")")
	})
	for {
		_, err := b.DownloadHTTP(url, filename, "http://127.0.0.1:1081")
		if !logs.PrintErr(err) {
			os.Remove(filename)
			break
		}
		<-time.After(time.Second * 5)
	}

	{
		logs.Debug("失败示例:")
		b = bar.New(100)
		go func(b *bar.Bar) {
			<-time.After(time.Second * 10)
			b.Close()
		}(b)
		b.Add(1).Flush()
		<-b.Done()
	}

	g.Input("请按回车键退出...")
}
