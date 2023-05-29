package main

import (
	"github.com/injoyai/base/oss"
	"github.com/injoyai/goutil/cmd/nac"
	"github.com/injoyai/goutil/string/bar"
	"github.com/injoyai/logs"
	"os"
	"path/filepath"
)

func main() {
	nac.Init()
	url := "https://github.com/injoyai/goutil/raw/main/cmd/in.exe"
	path, _ := os.Executable()
	filename := filepath.Join(filepath.Dir(path), "in.exe")
	for logs.PrintErr(bar.Download(url, filename)) {
	}
	oss.Input("升级成功,按任意键退出...")
}
