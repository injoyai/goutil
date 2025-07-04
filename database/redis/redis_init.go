package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/str"
	"github.com/redis/go-redis/v9"
	"time"
)

const Nil = redis.Nil

type (
	Config    = redis.Options
	StringCmd = redis.StringCmd
)

func New(addr, pwd string, db ...int) *Client {
	return NewClient(&Config{
		Addr:     addr,
		Password: pwd,
		DB:       conv.Default[int](0, db...),
	})
}

func WithCfg(key ...string) *Client {
	return WithDMap(cfg.GetDMap(conv.Default[string]("redis", key...)))
}

func WithDMap(m *conv.Map) *Client {
	return NewClient(&Config{
		Network:               m.GetString("network"),
		Addr:                  m.GetString("addr", fmt.Sprintf("%s:%s", m.GetString("host"), m.GetString("port"))),
		ClientName:            m.GetString("clientName"),
		Protocol:              m.GetInt("protocol"),
		Username:              m.GetString("username"),
		Password:              m.GetString("password", m.GetString("pwd")),
		DB:                    m.GetInt("db"),
		MaxRetries:            m.GetInt("maxRetries"),
		MinRetryBackoff:       m.GetDuration("minRetryBackoff"),
		MaxRetryBackoff:       m.GetDuration("maxRetryBackoff"),
		DialTimeout:           m.GetDuration("dialTimeout"),
		ReadTimeout:           m.GetDuration("readTimeout"),
		WriteTimeout:          m.GetDuration("writeTimeout"),
		ContextTimeoutEnabled: m.GetBool("contextTimeoutEnabled"),
		PoolFIFO:              m.GetBool("poolFIFO"),
		PoolSize:              m.GetInt("poolSize"),
		PoolTimeout:           m.GetDuration("poolTimeout"),
		MinIdleConns:          m.GetInt("minIdleConns"),
		MaxIdleConns:          m.GetInt("maxIdleConns"),
		MaxActiveConns:        m.GetInt("maxActiveConns"),
		ConnMaxIdleTime:       m.GetDuration("connMaxIdleTime"),
		ConnMaxLifetime:       m.GetDuration("connMaxLifetime"),
		DisableIndentity:      m.GetBool("disableIndentity"),
		IdentitySuffix:        m.GetString("identitySuffix"),
	})
}

func NewClient(op *Config) *Client {
	c := &Client{
		Client:          redis.NewClient(op),
		ctx:             context.Background(),
		CacheMap:        maps.NewSafe(),
		CacheExpiration: time.Second * 10,
	}
	c.Extend = conv.NewExtend(c)
	return c
}

type Client struct {
	*redis.Client
	ctx             context.Context
	CacheMap        *maps.Safe      //三级缓存,优先从内存,然后redis,然后调用函数
	CacheExpiration time.Duration   //缓存有效期,缓存失效的话,回去redis获取数据
	OnGetVarErr     func(err error) //获取var的错误信息,例改成panic,捕获到错误
	conv.Extend
}

func (this *Client) Ping() error {
	_, err := this.Client.Ping(this.ctx).Result()
	return err
}

func (this *Client) GetCmd(key string) *redis.StringCmd {
	return this.Client.Get(this.ctx, key)
}

func (this *Client) Get(key string) (string, error) {
	return this.GetCmd(key).Result()
}

// GetVar 实现接口,忽略了错误,并不安全
func (this *Client) GetVar(key string) *conv.Var {
	result := this.GetCmd(key)
	if result.Err() != nil {
		if result.Err() != Nil && this.OnGetVarErr != nil {
			this.OnGetVarErr(result.Err())
		}
		return conv.Nil()
	}
	return conv.New(this.GetCmd(key).Val())
}

func (this *Client) Set(key string, value any, expiration time.Duration) error {
	return this.Client.Set(this.ctx, key, value, expiration).Err()
}

// Cache 优先从内存中获取数据,不存在则尝试重redis中获取,小于等于0是不过期
func (this *Client) Cache(key string, fn func() (any, error), expiration time.Duration, cacheExpirations ...time.Duration) (any, error) {
	cacheExpiration := conv.Select[time.Duration](expiration > 0 && this.CacheExpiration > expiration, expiration, this.CacheExpiration)
	cacheExpiration = conv.Default[time.Duration](cacheExpiration, cacheExpirations...)
	return this.CacheMap.GetOrSetByHandler(key, func() (any, error) {
		s, err := this.Get(key)
		if err != nil && err.Error() != redis.Nil.Error() {
			return nil, err
		} else if err != nil {
			//假如redis中数据不存在,则使用函数生成数据,一般是从数据库获取
			data, err := fn()
			if err != nil {
				return nil, err
			}
			bs, err := json.Marshal(g.Map{"data": data})
			if err != nil {
				return nil, err
			}
			//保存数据到redis
			if err := this.Set(key, string(bs), expiration); err != nil {
				return nil, err
			}
			return data, nil
		}
		m := g.Map{}
		err = json.Unmarshal(str.Bytes(&s), &m)
		return m["data"], err
	}, cacheExpiration)
}
