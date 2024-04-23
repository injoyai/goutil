package runtime

import (
	"context"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	DefaultGoManage.SetLimit(10)
	for i := 0; i < 100; i++ {
		Try(func(ctx context.Context, args ...interface{}) {
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
