package runtimes

import (
	"github.com/injoyai/logs"
	"testing"
	"time"
)

func TestNewGoPool(t *testing.T) {
	p := NewGoPool(2)
	for i := 0; i < 10; i++ {
		p.Go(func() {
			logs.Debug(i)
			<-time.After(time.Second)
		})
	}
	t.Log(p.Current())
	<-time.After(time.Second)
	t.Log(p.Current())
	<-time.After(time.Second * 10)
	t.Log(p.Current())
}
