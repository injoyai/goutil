package runtimes

import (
	"context"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	DefaultGoManage.SetLimit(10)
	for i := 0; i < 100; i++ {
		Go(func(ctx context.Context, args ...interface{}) {
			//for {
			<-time.After(time.Second * 1)
			t.Log(args[0])
			//}

		}, i)
	}
	t.Log(DefaultGoManage.Len())
	<-time.After(time.Second * 10)
	t.Log(DefaultGoManage.Len())
}

func TestGoInfo(t *testing.T) {
	DefaultGoManage.SetLimit(10)
	f := func(ctx context.Context, args ...interface{}) {
		<-time.After(time.Second * 10)
	}
	item := Go(f)
	<-time.After(time.Second * 2)
	t.Log(item.StarTime)
	t.Log(item.FuncName())
	<-time.After(time.Second * 10)
	t.Log(item.StopTime)

}
