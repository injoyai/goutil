package in

import (
	"github.com/injoyai/base/maps"
	"sync"
	"time"
)

var (
	cache     *maps.Safe
	cacheOnce sync.Once
)

func CacheByHandler(key any, handler func() any, expiration time.Duration) any {
	cacheOnce.Do(func() {
		cache = maps.NewSafe()
	})
	value, err := cache.GetOrSetByHandler(key, func() (any, error) { return handler(), nil }, expiration)
	CheckErr(err)
	return value
}

func CacheDel(key ...any) {
	cacheOnce.Do(func() {
		cache = maps.NewSafe()
	})
	for _, v := range key {
		cache.Del(v)
	}
}
