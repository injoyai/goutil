package net_tools

import (
	"testing"
	"time"

	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/ios/v2"
	"github.com/injoyai/ios/v2/client"
	"github.com/injoyai/ios/v2/module/tcp"
)

func TestNewTCPClientEnable(t *testing.T) {
	e := NewTCPClientEnable(tcp.NewDial(":10086"), func(c *client.Client) {
		c.Logger.Enable()
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
