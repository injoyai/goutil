package main

import (
	"fmt"
	"github.com/injoyai/base/oss/shell"
	oss2 "github.com/injoyai/goutil/oss"
	"github.com/injoyai/logs"
	"github.com/spf13/cobra"
	"strings"
)

func handlerOpen(cmd *cobra.Command, args []string, flags *Flags) {
	if len(args) == 0 {
		fmt.Printf("请输入参数,例(in open hosts)")
		return
	}
	switch strings.ToLower(args[0]) {
	case "hosts":
		if shell.Start("C:\\Windows\\System32\\drivers\\etc\\hosts") != nil {
			shell.Start("C:\\Windows\\System32\\drivers\\etc\\")
		}
	case "injoy":
		shell.Start(oss2.UserDefaultDir())
	case "appdata":
		shell.Start(oss2.UserDataDir())
	case "startup":
		shell.Start(oss2.UserStartupDir())
	case "hfs", "downloader", "influxdb", "chromedriver":
		handlerInstall(cmd, args, flags)
		logs.PrintErr(shell.Start(args[0] + ".exe"))
	default:
		shell.Start(args[0])
	}
}
