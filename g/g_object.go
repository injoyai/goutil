package g

import (
	"context"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
	"math/rand"
	"sync"
	"time"
)

//========================================Rand========================================

var (
	ro sync.Once
	r  *rand.Rand
	rs = "abcdefghigklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ0123456789"
)

// Rand 随机数,单例模式,惰性加载
func Rand() *rand.Rand {
	ro.Do(func() {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	return r
}

// RandString 随机字符串
func RandString(length int, str ...string) string {
	xs := conv.Default[string](rs, str...)
	var s []byte
	for i := 0; i < length; i++ {
		n := Rand().Intn(len(xs))
		s = append(s, xs[n])
	}
	return string(s)
}

// RandInt 随机整数
func RandInt(min, max int) int {
	if max < min {
		return 0
	}
	return Rand().Intn(max-min) + min
}

// RandInt64 随机64位整数
func RandInt64(min, max int64) int64 {
	if max < min {
		return 0
	}
	return Rand().Int63n(max-min) + min
}

// RandFloat 随机浮点数
func RandFloat(min, max float64, d ...int) float64 {
	f := Rand().Float64()*(max-min) + min
	return Decimals(f, d...)
}

//========================================Context========================================

// Ctx context.Background
func Ctx(ctx ...context.Context) context.Context {
	if len(ctx) > 0 && ctx[0] != nil {
		return ctx[0]
	}
	return context.Background()
}

// WithCancel context.WithCancel
func WithCancel(ctx ...context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(Ctx(ctx...))
}

// WithTimeout context.WithTimeout
func WithTimeout(timeout time.Duration, ctx ...context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(Ctx(ctx...), timeout)
}

func WithDeadline(deadline time.Time, ctx ...context.Context) (context.Context, context.CancelFunc) {
	return context.WithDeadline(Ctx(ctx...), deadline)
}

func WithValue(key, val any, ctx ...context.Context) context.Context {
	return context.WithValue(Ctx(ctx...), key, val)
}

//========================================Other========================================

// Closer 安全的关闭,原子操作,实现接口
func Closer() *safe.Closer { return safe.NewCloser() }

// Runner 安全的执行,避免重复执行,实现接口
func Runner(fn func(ctx context.Context) error) *safe.Runner { return safe.NewRunner(fn) }
