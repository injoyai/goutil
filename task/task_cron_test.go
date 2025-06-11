package task

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	x := New().Start()
	x.SetTask("", "20,50 54 15 * * *", func() {
		t.Log(time.Now().Format("2006-01-02 15:04:05: test"))
	})
	select {}
}
