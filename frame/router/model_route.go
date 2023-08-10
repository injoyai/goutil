package router

import (
	"sync"
)

func newRoute() *Route {
	return &Route{
		m: make(map[string]Handler),
	}
}

type Route struct {
	all Handler
	m   map[string]Handler
	mu  sync.RWMutex
}

func (r *Route) Clone() map[string]Handler {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m := make(map[string]Handler)
	if _, ok := r.m["ALL"]; ok {
		return map[string]Handler{"ALL": r.all}
	}
	for i, v := range r.m {
		m[i] = v
	}
	return m
}

func (r *Route) IsEmpty() bool {
	return r.all == nil && r.Len() == 0
}

func (r *Route) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.m)
}

func (r *Route) Get(method string) Handler {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.all != nil {
		return r.all
	}
	return r.m[method]
}

func (r *Route) Set(method string, handlerFunc Handler) *Route {
	if handlerFunc != nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		if method == "ALL" {
			r.all = handlerFunc
		}
		r.m[method] = handlerFunc
	}
	return r
}
