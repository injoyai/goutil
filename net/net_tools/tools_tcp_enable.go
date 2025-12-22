package net_tools

import (
	"context"

	"github.com/injoyai/base/safe"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
)

func NewTCPClientEnable(dial ios.DialFunc, options ...client.Option) *TCPClientEnable {
	return NewTCPClientEnableWithContext(context.Background(), dial, options...)
}

// NewTCPClientEnableWithContext 单例的TCP客户端启用禁用
func NewTCPClientEnableWithContext(ctx context.Context, dial ios.DialFunc, options ...client.Option) *TCPClientEnable {
	return &TCPClientEnable{
		DialFunc: dial,
		Options:  options,
		ctx:      ctx,
		Runner:   safe.NewRunnerWithContext(ctx, nil),
	}
}

type TCPClientEnable struct {
	ios.DialFunc
	Options []client.Option
	ctx     context.Context
	*safe.Runner
}

func (this *TCPClientEnable) Enable() error {

	this.Runner.SetFunc(func(ctx context.Context) error {
		return client.RedialContext(this.ctx, this.DialFunc, func(c *client.Client) {
			go func() {
				select {
				case <-ctx.Done():
					c.CloseAll()
				case <-c.Done():
				}
			}()
			c.SetOption(this.Options...)
		}).Run(context.Background())
	})
	go this.Runner.Run()
	return nil
}

func (this *TCPClientEnable) Disable() {
	this.Runner.Stop(true)
}
