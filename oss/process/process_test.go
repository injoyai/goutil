package process

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/injoyai/ios/v2"
)

func TestNew(t *testing.T) {
	p := NewPowershell("ping -t 8.8.8.8")
	p.SetStdout(ios.WriteFunc(func(p []byte) (int, error) {
		return fmt.Println(string(p))
	}))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	t.Error(p.Run(ctx))
}
