package g

import (
	"context"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/conv/codec"
	"math/rand"
	"time"
)

//========================================Rand========================================

var (
	r  *rand.Rand
	rs = "abcdefghigklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ0123456789"
)

// Rand 随机数,单例模式,惰性加载
func Rand() *rand.Rand {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return r
}

// RandString 随机字符串
func RandString(length int, str ...string) string {
	r := Rand()
	xs := conv.GetDefaultString(rs, str...)
	var s []byte
	for i := 0; i < length; i++ {
		n := r.Intn(len(xs))
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

func WithValue(key, val interface{}, ctx ...context.Context) context.Context {
	return context.WithValue(Ctx(ctx...), key, val)
}

//========================================Time========================================

// Now 当前时间
func Now() time.Time { return time.Now() }

// Date 当前日期
func Date() (year int, month time.Month, day int) { return time.Now().Date() }

// Unix 当前时间戳
func Unix() int64 { return time.Now().Unix() }

// UnixNano 当前纳秒
func UnixNano() int64 { return time.Now().UnixNano() }

// Year 年[1970-]
func Year() int { return time.Now().Year() }

// Month 月[1-12]
func Month() int { return int(time.Now().Month()) }

// Day 日[1-31]
func Day() int { return time.Now().Day() }

// Hour 时[0-23]
func Hour() int { return time.Now().Hour() }

// Minute 分[0-59]
func Minute() int { return time.Now().Minute() }

// Second 秒[0-59]
func Second() int { return time.Now().Second() }

//========================================Other========================================

// Cfg 读取配置文件
func Cfg(path string, codec ...codec.Interface) *cfg.Entity { return cfg.New(path, codec...) }

// Chan chans.NewEntity
func Chan(num int, cap ...int) *chans.Entity { return chans.NewEntity(num, cap...) }

// NewCloser 安全的关闭,原子操作
func NewCloser() *safe.Closer { return safe.NewCloser() }
