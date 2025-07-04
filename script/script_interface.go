package script

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
)

type Func func(*Args) (any, error)
type Option func(c Client)

type Client interface {
	Exec(text string, option ...func(i Client)) (any, error)
	Set(key string, value any) error
	SetFunc(key string, value Func) error
	Close() error
	Tag() *maps.Safe
}

type Args struct {
	This Client
	Args []*conv.Var
}

func (this *Args) Len() int {
	return len(this.Args)
}

func (this *Args) Get(idx int) *conv.Var {
	if this.Len() > idx-1 && idx > 0 {
		return this.Args[idx-1]
	}
	return conv.Nil()
}

func (this *Args) GetInt(idx int, def ...int) int {
	return this.Get(idx).Int(def...)
}

func (this *Args) GetInt64(idx int, def ...int64) int64 {
	return this.Get(idx).Int64(def...)
}

func (this *Args) GetFloat64(idx int, def ...float64) float64 {
	return this.Get(idx).Float64(def...)
}

func (this *Args) GetString(idx int, def ...string) string {
	return this.Get(idx).String(def...)
}

func (this *Args) GetBytes(idx int, def ...[]byte) []byte {
	return this.Get(idx).Bytes(def...)
}

func (this *Args) GetBool(idx int, def ...bool) bool {
	return this.Get(idx).Bool(def...)
}

func (this *Args) GetMap(idx int, def ...map[string]any) map[string]any {
	return this.Get(idx).Map(def...)
}

func (this *Args) GetGMap(idx int, def ...map[string]any) map[string]any {
	return this.Get(idx).GMap(def...)
}

func (this *Args) GetDMap(idx int, def ...any) *conv.Map {
	return this.Get(idx).DMap(def...)
}

func (this *Args) GetArray(idx int, def ...[]any) []any {
	return this.Get(idx).Interfaces(def...)
}

func (this *Args) Interfaces() []any {
	list := make([]any, 0)
	for _, v := range this.Args {
		list = append(list, v.Val())
	}
	return list
}
