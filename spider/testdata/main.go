package main

import (
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/spider"
	"github.com/injoyai/logs"
)

func main() {
	logs.PrintErr(spider.New(
		oss.UserInjoyDir("/downloader/browser/chrome/chrome.exe"),
		oss.UserInjoyDir("/downloader/browser/chrome/chromedriver.exe"),
	).Run(func(w *spider.WebDriver) error {
		w.Open("http://www.baidu.com")
		w.WaitMin(2)
		return nil
	}))
}
