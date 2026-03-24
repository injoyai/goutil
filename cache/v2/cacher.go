package cache

import (
	"context"
	"time"

	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/redis/go-redis/v9"
)

/*

试用

*/

type Cacher[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, value V, expiration ...time.Duration) error
	Del(key K) error
}

func NewRedis(client *redis.Client) Cacher[string, any] {
	return &_redis{Client: client}
}

func NewMemory[K comparable, V any]() Cacher[K, V] {
	return &_memory[K, V]{Generic: maps.NewGeneric[K, V]()}
}

func NewDisk(name string, groups ...string) Cacher[string, any] {
	return &_disk{File: NewFile(name, groups...)}
}

/*



 */

type _redis struct {
	*redis.Client
}

func (this *_redis) Get(key string) (any, error) {
	return this.Client.Get(context.Background(), conv.String(key)).Result()
}

func (this *_redis) Set(key string, value any, expiration ...time.Duration) error {
	exp := conv.Default(-1, expiration...)
	return this.Client.Set(context.Background(), key, value, exp).Err()
}

func (this *_redis) Del(key string) error {
	return this.Client.Del(context.Background(), key).Err()
}

/**/

type _memory[K comparable, V any] struct {
	*maps.Generic[K, V]
}

func (this *_memory[K, V]) Get(key K) (V, error) {
	val, _ := this.Generic.Get(key)
	return val, nil
}

func (this *_memory[K, V]) Set(key K, value V, expiration ...time.Duration) error {
	this.Generic.Set(key, value, expiration...)
	return nil
}

func (this *_memory[K, V]) Del(key K) error {
	this.Generic.Del(key)
	return nil
}

/**/

type _disk struct {
	*File
}

func (this *_disk) Get(key string) (any, error) {
	return this.File.GetInterface(key), nil
}

func (this *_disk) Set(key string, value any, expiration ...time.Duration) error {
	this.File.Set(key, value)
	return nil
}

func (this *_disk) Del(key string) error {
	this.File.Map.Del(key)
	return nil
}
