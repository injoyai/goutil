package ctx

import (
	"context"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/base/safe"
	"time"
)

var _ context.Context = (*Ctx)(nil)

func New() *Ctx {
	return &Ctx{
		Safe:   maps.NewSafe(),
		Closer: safe.NewCloser(),
	}
}

type Ctx struct {
	*maps.Safe
	*safe.Closer
}

func (this *Ctx) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (this *Ctx) Value(key any) any {
	return this.Safe.MustGet(key)
}
