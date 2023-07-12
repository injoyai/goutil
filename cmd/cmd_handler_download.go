package main

import (
	"context"
	"fmt"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/oss"
	"github.com/injoyai/goutil/cache"
	oss2 "github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/protocol/m3u8"
	"github.com/injoyai/goutil/string/bar"
	"github.com/injoyai/logs"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func handlerDownload(cmd *cobra.Command, args []string, flags *Flags) {
	if len(args) == 0 {
		fmt.Println("请输入下载的内容")
		return
	}
	resource := args[0]
	filename := fmt.Sprintf("./%s.exe", strings.ToLower(resource))
	if len(resource) == 0 {
		fmt.Println("请输入下载的内容")
		return
	}
	switch strings.ToLower(resource) {

	case "in":

		url := "https://github.com/injoyai/goutil/raw/main/cmd/in.exe"
		logs.PrintErr(bar.Download(url, filename))

	case "upgrade":

		filename = "./in_upgrade.exe"
		logs.PrintErr(oss.New(filename, upgrade))

	case "upx":

		logs.PrintErr(oss.New(filename, upx))

	case "rsrc":

		logs.PrintErr(oss.New(filename, rsrc))

	case "chromedriver":

		if _, err := installChromedriver(oss2.UserDefaultDir(), flags.GetBool("download")); err != nil {
			log.Printf("[错误] %s", err.Error())
		}

	case "downloader":

		url := "https://github.com/injoyai/downloader/releases/latest/download/downloader.exe"
		logs.PrintErr(bar.Download(url, filename))

	case "swag":

		logs.PrintErr(oss.New(filename, swag))

	case "hfs":

		logs.PrintErr(oss.New(filename, hfs))

	case "influxdb":

		url := "https://dl.influxdata.com/influxdb/releases/influxdb2-2.7.1-windows-amd64.zip"
		logs.PrintErr(bar.Download(url, filename+".zip"))
		logs.PrintErr(DecodeZIP(filename+".zip", "./"))
		os.Remove(filename + ".zip")

	case "mingw64":

		//https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z

	default:

		switch true {
		case strings.Contains(resource, ".m3u8"),
			strings.Contains(resource, ".m3u8?"):

			cache.DefaultDir = oss2.UserDefaultDir() + "/config/"
			cfg := cache.NewFile("cmd_download", "cmd").Sync()
			dir := flags.GetString("dir", cfg.GetString("dir"))
			cfg.Set("dir", dir)
			os.MkdirAll(dir, 0666)
			filename = filepath.Join(dir, flags.GetString("output", filepath.Base(resource)))
			goroute := flags.GetInt("goroute")
			tryNum := flags.GetInt("try")

			list, err := m3u8.New(resource)
			if err != nil {
				log.Printf("[错误] %v", err)
				return
			}
			input := []_downloadRun(nil)
			for _, v := range list {
				input = append(input, v)
			}
			result := newDownload(goroute, tryNum, input)
			bytes := []byte(nil)
			for _, v := range result {
				bytes = append(bytes, v...)
			}
			logs.PrintErr(oss.New(filename, bytes))

		default:

			logs.PrintErr(bar.Download(resource, filename))

		}

	}
}

func newDownload(num, tryNum int, runs []_downloadRun) [][]byte {
	b := bar.New(int64(len(runs)))
	cache := make([][]byte, len(runs))
	queue := chans.NewEntity(num, len(runs))
	queue.SetHandler(func(ctx context.Context, no, count int, data interface{}) {
		x := data.(*_downloadItem)
		var bytes []byte
		var err error
		var retry int
		for bytes, err = x.run(); err != nil && retry < tryNum; retry++ {
			<-time.After(time.Second)
		}
		b.Add(1)
		cache[x.idx] = bytes
	})
	go func() {
		for idx, run := range runs {
			queue.Do(&_downloadItem{
				idx: idx,
				run: run.GetBytes,
			})
		}
	}()
	<-b.Run()
	return cache
}

type _downloadRun interface {
	GetBytes() ([]byte, error)
}

type _downloadItem struct {
	idx int
	run func() ([]byte, error)
}
