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

// CacheByHandler 尝试从缓存中获取数据,如果不存在则通过函数获取,执行函数时,其他相同的key会等待此次结果
func CacheByHandler(key any, handler func() any, expiration time.Duration) any {
	requestCacheOnce.Do(func() {
		RequestCache = maps.NewSafe()
	})
	value, err := RequestCache.GetOrSetByHandler(key, func() (any, error) { return handler(), nil }, expiration)
	CheckErr(err)
	return value
}

// CacheDel 删除缓存数据
func CacheDel(key ...any) {
	requestCacheOnce.Do(func() {
		RequestCache = maps.NewSafe()
	})
	for _, v := range key {
		RequestCache.Del(v)
	}
}

// CacheSet 设置缓存,覆盖缓存
func CacheSet(key any, value any, expiration time.Duration) {
	requestCacheOnce.Do(func() {
		RequestCache = maps.NewSafe()
	})
	RequestCache.Set(key, value, expiration)
}
