package bar

import "github.com/injoyai/base/chans"

func NewCoroutine(total, limit int, op ...Option) Coroutine {
	b := New(op...)
	b.SetTotal(int64(total))
	b.Flush()
	return &coroutine{
		Bar: b,
		wg:  chans.NewWaitLimit(limit),
	}
}

type coroutine struct {
	Bar
	wg chans.WaitLimit
}

func (this *coroutine) Wait() {
	this.wg.Wait()
}

func (this *coroutine) Go(f func()) {
	this.GoRetry(func() error {
		f()
		return nil
	}, 1)
}

func (this *coroutine) GoRetry(f func() error, retry int) {
	if f == nil {
		return
	}
	this.wg.Add()
	go func() {
		defer func() {
			this.Bar.Add(1)
			this.Bar.Flush()
			this.wg.Done()
		}()
		for i := 0; i < retry; i++ {
			if err := f(); err == nil {
				return
			}
		}
	}()
}
