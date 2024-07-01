package g

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	Retry(func() error {
		t.Log(time.Now().Format("15:04:05"), "重试")
		return errors.New("")
	}, -1, WithRetreat32)
}
