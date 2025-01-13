package cache

import (
	"context"
	"github.com/injoyai/base/maps"
	"github.com/redis/go-redis/v9"
	"time"
)

/*

试用

*/

type Interface interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiration ...time.Duration) error
	Del(key string) error
}

func NewRedisCacher(client *redis.Client) Interface {
	return &_redis{Client: client}
}

func NewMapCacher(m ...maps.Map) Interface {
	return &_map{Safe: maps.NewSafe(m...)}
}

func NewFileCacher(name string, groups ...string) Interface {
	return &_file{File: newFile(name, groups...)}
}

/*



 */

type _redis struct {
	*redis.Client
}

func (this *_redis) Get(key string) (interface{}, error) {
	s, err := this.Client.Get(context.Background(), key).Result()
	return s, err
}

func (this *_redis) Set(key string, value interface{}, expiration ...time.Duration) error {
	if len(expiration) == 0 {
		return this.Client.Set(context.Background(), key, value, -1).Err()
	}
	return this.Client.Set(context.Background(), key, value, expiration[0]).Err()
}

func (this *_redis) Del(key string) error {
	return this.Client.Del(context.Background(), key).Err()
}

/**/

type _map struct {
	*maps.Safe
}

func (this *_map) Get(key string) (interface{}, error) {
	val, _ := this.Safe.Get(key)
	return val, nil
}

func (this *_map) Set(key string, value interface{}, expiration ...time.Duration) error {
	this.Safe.Set(key, value, expiration...)
	return nil
}

func (this *_map) Del(key string) error {
	this.Safe.Del(key)
	return nil
}

/**/

type _file struct {
	*File
}

func (this *_file) Get(key string) (interface{}, error) {
	return this.File.GetInterface(key), nil
}

func (this *_file) Set(key string, value interface{}, expiration ...time.Duration) error {
	this.File.Set(key, value)
	return nil
}

func (this *_file) Del(key string) error {
	this.File.Map.Del(key)
	return nil
}
