package net_tools

import (
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial"
	"testing"
	"time"
)

func TestNewTCPClientEnable(t *testing.T) {
	e := NewTCPClientEnable(dial.WithTCP(":10086"), func(c *io.Client) {
		c.Debug()
		c.GoTimerWriteString(time.Second, time.Now().Format("15:04:05"))
	})
	<-time.After(time.Second * 5)
	t.Log("启用")
	t.Log(e.Enable())
	<-time.After(time.Second * 10)
	t.Log("禁用")
	e.Disable()
	e.Enable()
	oss.Wait()
}
