package router

import "sync"

type Bind struct {
	m  map[int]Handler
	mu sync.RWMutex
}

func newBind() *Bind {
	return &Bind{
		m: make(map[int]Handler),
	}
}

func (b *Bind) Set(code int, fn Handler) *Bind {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.m[code] = fn
	return b
}

func (b *Bind) Get(code int) Handler {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.m[code]
}

func (b *Bind) GetAndDo(code int, r *Request) {
	if handler := b.Get(code); handler != nil {
		handler(r)
	}
}
