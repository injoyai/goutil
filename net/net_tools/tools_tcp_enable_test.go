package net_tools

import (
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/client/dial"
	"testing"
	"time"
)

func TestNewTCPClientEnable(t *testing.T) {
	e := NewTCPClientEnable(dial.WithTCP(":10086"), func(c *client.Client) {
		c.Logger.Debug()
		c.GoTimerWriter(time.Second, func(w ios.MoreWriter) error {
			_, err := w.WriteString(time.Now().Format("15:04:05"))
			return err
		})
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
