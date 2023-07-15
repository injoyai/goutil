package main

import (
	"github.com/injoyai/goutil/oss"
	"log"
)

func main() {
	oss.ListenExit(func() { log.Println(1) })
	oss.ListenExit(func() { log.Println(2) })
	oss.Wait()
}
