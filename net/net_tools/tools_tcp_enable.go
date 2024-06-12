package net_tools

import (
	"context"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/io"
)

func NewTCPClientEnable(dial io.DialFunc, options ...io.OptionClient) *TCPClientEnable {
	return NewTCPClientEnableWithContext(context.Background(), dial, options...)
}

// NewTCPClientEnableWithContext 单例的TCP客户端启用禁用
func NewTCPClientEnableWithContext(ctx context.Context, dial io.DialFunc, options ...io.OptionClient) *TCPClientEnable {
	return &TCPClientEnable{
		DialFunc: dial,
		Options:  options,
		ctx:      ctx,
		Runner:   safe.NewRunnerWithContext(ctx, nil),
	}
}

type TCPClientEnable struct {
	io.DialFunc
	Options []io.OptionClient
	ctx     context.Context
	*safe.Runner
}

func (this *TCPClientEnable) Enable() error {

	this.Runner.SetFunc(func(ctx context.Context) error {
		return io.RedialWithContext(this.ctx, this.DialFunc, func(c *io.Client) {
			go func() {
				select {
				case <-ctx.Done():
					c.CloseAll()
				case <-c.Done():
				}
			}()
			c.SetOptions(this.Options...)
		}).Run()
	})
	go this.Runner.Run()
	return nil
}

func (this *TCPClientEnable) Disable() {
	this.Runner.Stop(true)
}
