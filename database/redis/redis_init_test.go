package redis

import (
	"context"
	"testing"
	"time"
)

func TestClient_Cache(t *testing.T) {
	c := New("127.0.0.1:6379", "")
	c.Del(context.Background(), "key")
	for i := 0; i < 10; i++ {
		val, err := c.Cache("key", func() (any, error) {
			t.Log("handler generate data")
			return map[string]any{
				"time": time.Now().String(),
			}, nil
		}, time.Second)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(val)
		<-time.After(time.Second * 6)
	}
}

func TestWithCfg(t *testing.T) {
	t.Log(WithCfg().Ping())
}
