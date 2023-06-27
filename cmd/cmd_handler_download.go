package main

import (
	"context"
	"fmt"
	"github.com/injoyai/base/chans"
	oss2 "github.com/injoyai/base/oss"
	"github.com/injoyai/goutil/cache"
	"github.com/injoyai/goutil/oss"
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
		fmt.Println("请输入下载内容")
	}

	cache.DefaultDir = oss.UserDefaultDir() + "/config/"
	cfg := cache.NewFile("cmd_download", "cmd").Sync()
	dir := flags.GetString("dir", cfg.GetString("dir"))
	cfg.Set("dir", dir)
	os.MkdirAll(dir, 0666)
	filename := filepath.Join(dir, flags.GetString("output", filepath.Base(args[0])))
	goroute := flags.GetInt("goroute")
	tryNum := flags.GetInt("try")

	switch true {
	case strings.Contains(args[0], ".m3u8"),
		strings.Contains(args[0], ".m3u8?"):

		list, err := m3u8.New(args[0])
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
		logs.PrintErr(oss2.New(filename, bytes))

	default:

		logs.PrintErr(bar.Download(args[0], filename))

	}

}

func newDownload(num, tryNum int, runs []_downloadRun) [][]byte {
	b := bar.New().SetTotalSize(float64(len(runs)))
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
