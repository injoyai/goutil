package main

import (
	"github.com/injoyai/goutil/other/keyboard"
	"github.com/injoyai/logs"
)

func main() {

	//c := keyboard.ListenKey(keyboard.KeySpace)
	//for range c {
	//	logs.Debug()
	//}

	logs.Debug(keyboard.ListenFunc(func(key keyboard.KeyEvent) {

		logs.Debug(key.Key, "  ", string([]rune{key.Rune}))
	}))
}
