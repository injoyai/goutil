package main

import (
	"github.com/injoyai/goutil/string/bar"
	"github.com/injoyai/logs"
)

func init() {
	logs.DefaultErr.SetWriter(logs.Stdout)
}

func main() {
	bar.Demo()
}
