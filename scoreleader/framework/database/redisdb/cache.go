package redisdb

// import (
// 	"encoding/json"
// 	"errors"
// 	"sync"
// 	"time"

// 	redis "github.com/redis/go-redis/v9"
// 	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
// )

// const (
// 	ErrRedisConfigInvalid  = "RedisConfig is invalid"
// 	ErrCacheNotInitialized = "Cache is not initialized"
// 	ErrRedisConfigIsNil    = "RedisConfig is nil"
// 	ErrRedisIsNil          = "Redis is nil"
// 	ErrNotCached           = "Not cached"
// )

// type Cache struct {
// 	enable     bool
// 	db         *Redis
// 	err        error
// 	timeoutSec time.Duration
// }

// // RedisConfig Cache 설정
// type RedisConfig struct {
// 	Host       string `json:"host" yaml:"host"`
// 	Port       int    `json:"port" yaml:"port"`
// 	Password   string `json:"password" yaml:"password"`
// 	DefaultDB  int    `json:"database" yaml:"database"`
// 	PoolSize   int    `json:"poolsize" yaml:"poolsize"`
// 	TimeoutSec int    `json:"timeout_sec" yaml:"timeout_sec"`
// }

// var gCache *Cache
// var gOnceCache sync.Once

// func GetCache(options ...*RedisConfig) *Cache {
// 	gOnceCache.Do(func() {
// 		gCache = new(Cache)

// 		if len(options) != 1 {
// 			gCache.err = errors.New(ErrRedisConfigInvalid)
// 			return
// 		}

// 		gCache = NewCache(options[0])
// 	})

// 	return gCache
// }

// func NewCache(option *RedisConfig) *Cache {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Error("Recovered panic", r)

// 		}
// 	}()

// 	cache := new(Cache)

// 	if option == nil {
// 		cache.err = errors.New(ErrRedisConfigIsNil)
// 	}
// 	cache.db = NewRedis(option.Host, option.Port, option.Password, option.DefaultDB, option.PoolSize)
// 	cache.timeoutSec = time.Duration(option.TimeoutSec) * time.Second
// 	cache.enable, cache.err = cache.db.Ping()

// 	return cache
// }

// func (o Cache) Enable() bool {
// 	return o.enable
// }

// func (o *Cache) Close() error {
// 	if !o.enable && o.db == nil {
// 		return nil
// 	}

// 	if o.db == nil {
// 		return errors.New(ErrRedisIsNil)
// 	}
// 	o.err = o.db.Close()
// 	o.enable = false
// 	return o.err
// }

// func (o *Cache) Del(key string) error {
// 	if !o.enable || o.db == nil {
// 		return errors.New(ErrCacheNotInitialized)
// 	}

// 	if _, err := o.db.Del(key); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // Set 만료시간을 갖는 데이터 저장
// // expiration이 지정되지 않을 경우 cache에 지정된 기본 만료시간(o.timeoutSec)을 설정한다.
// func (o *Cache) Set(key string, data interface{}, expiration ...time.Duration) error {
// 	if !o.enable || o.db == nil {
// 		return errors.New(ErrCacheNotInitialized)
// 	}

// 	expirationSec := o.timeoutSec

// 	if len(expiration) == 1 {
// 		expirationSec = expiration[0]
// 	}

// 	buff, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}

// 	if _, err := o.db.SetAndExpire(key, string(buff[:]), expirationSec); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (o *Cache) SetNX(key string, value string, expiration time.Duration) (bool, error) {
// 	if !o.enable || o.db == nil {
// 		return false, errors.New(ErrCacheNotInitialized)
// 	}

// 	return o.db.SetNX(key, value, expiration)
// }

// func (o *Cache) Get(key string, data interface{}) error {
// 	if !o.enable || o.db == nil {
// 		return errors.New(ErrCacheNotInitialized)
// 	}

// 	str, err := o.db.Get(key)
// 	if err != nil {
// 		if IsNil(err) {
// 			return errors.New(ErrNotCached)
// 		}
// 		return err
// 	}

// 	if err := json.Unmarshal([]byte(str), &data); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (o *Cache) HGet(key string, field string, data interface{}) error {
// 	if !o.enable || o.db == nil {
// 		return errors.New(ErrCacheNotInitialized)
// 	}

// 	str, err := o.db.HGet(key, field)
// 	if err != nil {
// 		if IsNil(err) {
// 			return errors.New(ErrNotCached)
// 		}
// 		return err
// 	}

// 	if err := json.Unmarshal([]byte(str), &data); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (o *Cache) HMGet(key string, fields []string, data []interface{}) error {
// 	if !o.enable || o.db == nil {
// 		return errors.New(ErrCacheNotInitialized)
// 	}

// 	vlist, err := o.db.HMGet(key, fields...)
// 	if err != nil {
// 		if IsNil(err) {
// 			return errors.New(ErrNotCached)
// 		}
// 		return err
// 	}

// 	for i, v := range vlist {
// 		if err := json.Unmarshal([]byte(v.(string)), &data[i]); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (o *Cache) HSet(key string, field string, value interface{}) error {
// 	if !o.enable || o.db == nil {
// 		return errors.New(ErrCacheNotInitialized)
// 	}

// 	buff, err := json.Marshal(value)
// 	if err != nil {
// 		return err
// 	}

// 	if err := o.db.HSet(key, field, string(buff[:])); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (o *Cache) HMSet(key string, fields map[string]interface{}) error {
// 	if !o.enable || o.db == nil {
// 		return errors.New(ErrCacheNotInitialized)
// 	}

// 	m := make(map[string]string)

// 	for field, value := range fields {
// 		buff, err := json.Marshal(value)
// 		if err != nil {
// 			return err
// 		}
// 		m[field] = string(buff[:])
// 	}

// 	if err := o.db.HMSet(key, m); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (o *Cache) HKeys(key string) ([]string, error) {
// 	if !o.enable || o.db == nil {
// 		return nil, errors.New(ErrCacheNotInitialized)
// 	}

// 	return o.db.HKeys(key)
// }

// // HScan 은 value 의 타입을 특정할 수 없기 때문에 다른 Cache 메서드들과 다르게 Unmarshal 하지 않은 결과를 리턴한다
// func (o *Cache) HScan(key string, cursor uint64, match string, count int64) (map[string]string, error) {
// 	if !o.enable || o.db == nil {
// 		return nil, errors.New(ErrCacheNotInitialized)
// 	}

// 	m := make(map[string]string)
// 	list, err := o.db.HScan(key, cursor, match, count)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for idx, val := range list {
// 		if idx%2 == 1 {
// 			m[list[idx-1]] = val
// 		}
// 	}

// 	return m, err
// }

// func (o *Cache) HDel(key string, field string) error {
// 	if !o.enable || o.db == nil {
// 		return errors.New(ErrCacheNotInitialized)
// 	}

// 	return o.db.HDel(key, field)
// }

// func (o *Cache) HIncrBy(key, field string, incr int64) (int64, error) {
// 	if !o.enable || o.db == nil {
// 		return 0, errors.New(ErrCacheNotInitialized)
// 	}

// 	return o.db.HIncrBy(key, field, incr)
// }

// func (o *Cache) Truncate(key string) error {
// 	if !o.enable || o.db == nil {
// 		return errors.New(ErrCacheNotInitialized)
// 	}

// 	itemList, err := o.db.Keys(key)
// 	if err != nil {
// 		return err
// 	}

// 	for _, item := range itemList {
// 		if _, err := o.db.Del(item); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (o *Cache) Pipeline() *redis.Pipeline {
// 	if !o.enable || o.db == nil {
// 		return nil
// 	}

// 	return o.Pipeline()
// }
