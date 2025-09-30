package main

import (
	"context"
	"errors"
	"github.com/injoyai/goutil/str/bar/v2"
	"github.com/injoyai/goutil/task"
	"time"
)

func main() {
	r := task.NewRange[any]()
	r.SetBar(bar.New())
	r.SetCoroutine(10)
	r.SetRetryInterval(func(i int) time.Duration {
		return time.Duration(i) * time.Second
	})
	r.Append(func(ctx context.Context) (any, error) {
		<-time.After(time.Second)
		return nil, errors.New("error message")
	})
	r.Append(func(ctx context.Context) (any, error) {
		<-time.After(time.Second)
		return nil, nil
	})
	r.Run()
}
