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

func CacheByHandler(key interface{}, handler func() interface{}, expiration time.Duration) interface{} {
	cacheOnce.Do(func() {
		cache = maps.NewSafe()
	})
	value, err := cache.GetOrSetByHandler(key, func() (interface{}, error) { return handler(), nil }, expiration)
	CheckErr(err)
	return value
}

func CacheDel(key ...interface{}) {
	cacheOnce.Do(func() {
		cache = maps.NewSafe()
	})
	for _, v := range key {
		cache.Del(v)
	}
}
