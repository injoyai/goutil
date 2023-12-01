package in

import (
	"github.com/injoyai/base/maps"
	"sync"
	"time"
)

var (
	RequestCache     *maps.Safe
	requestCacheOnce sync.Once
)

func CacheByHandler(key interface{}, handler func() interface{}, expiration time.Duration) interface{} {
	requestCacheOnce.Do(func() {
		RequestCache = maps.NewSafe()
	})
	value, err := RequestCache.GetOrSetByHandler(key, func() (interface{}, error) { return handler(), nil }, expiration)
	CheckErr(err)
	return value
}

func CacheDel(key ...interface{}) {
	requestCacheOnce.Do(func() {
		RequestCache = maps.NewSafe()
	})
	for _, v := range key {
		RequestCache.Del(v)
	}
}
