package chrome

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/injoyai/goutil/oss"
	"time"
)

type (
	ContextOption = chromedp.ContextOption
	QueryOption   = chromedp.QueryOption
	Action        = chromedp.Action
	ActionFunc    = chromedp.ActionFunc
)

func New(opts ...ContextOption) *Chrome {
	return WithContext(context.Background(), opts...)
}

func WithContext(ctx context.Context, opts ...ContextOption) *Chrome {
	ctx, cancel := chromedp.NewContext(ctx, opts...)
	return &Chrome{
		ctx:    ctx,
		cancel: cancel,
	}
}

type Chrome struct {
	ctx    context.Context
	cancel context.CancelFunc
	action []chromedp.Action
}

func (this *Chrome) Run() error {
	return chromedp.Run(this.ctx, this.action...)
}

func (this *Chrome) Action(action ...Action) *Chrome {
	this.action = append(this.action, action...)
	return this
}

func (this *Chrome) Get(url string) *Chrome {
	return this.Action(chromedp.Navigate(url))
}

func (this *Chrome) Click(sel any, opts ...QueryOption) *Chrome {
	return this.Action(chromedp.Click(sel, opts...))
}

func (this *Chrome) Sleep(d time.Duration) *Chrome {
	return this.Action(chromedp.Sleep(d))
}

func (this *Chrome) SaveScreenshot(sel any, filename string) *Chrome {
	return this.Action(chromedp.ActionFunc(func(ctx context.Context) error {
		var buf []byte
		if err := chromedp.Screenshot(sel, &buf, chromedp.NodeVisible).Do(ctx); err != nil {
			return err
		}
		return oss.New(filename, buf)
	}))
}
