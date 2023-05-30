package in

import (
	"github.com/injoyai/base/maps"
	"time"
)

var cache *maps.Safe

func CacheByHandler(key interface{}, handler func() interface{}, expiration time.Duration) interface{} {
	if cache == nil {
		cache = maps.NewSafe()
	}
	value, err := cache.GetOrSetByHandler(key, func() (interface{}, error) { return handler(), nil }, expiration)
	CheckErr(err)
	return value
}

func CacheDel(key ...interface{}) {
	for _, v := range key {
		cache.Del(v)
	}
}
