package redisdb

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"stove-gitlab.sginfra.net/backend/template/framework/util/log"

	redis "github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
)

type PubSubMessage *redis.Message

// RedisConfig Cache 설정
type RedisConfig struct {
	Host       string `json:"host" yaml:"host"`
	Port       int    `json:"port" yaml:"port"`
	Password   string `json:"password" yaml:"password"`
	DefaultDB  int    `json:"database" yaml:"database"`
	PoolSize   int    `json:"poolsize" yaml:"poolsize"`
	TimeoutSec int    `json:"timeout_sec" yaml:"timeout_sec"`
}

type Redis struct {
	client    *redis.Client
	pubsub    *redis.PubSub
	receiveCh chan PubSubMessage
}

var gRedis *Redis
var gOnceRedis sync.Once

func GetRedis(params ...string) (*Redis, error) {
	var err error
	gOnceRedis.Do(func() {
		host := ""
		port := int(6379)
		password := ""
		defaultDB := 0
		poolSize := 10

		if len(params) >= 1 {
			host = params[0]
		}

		if len(params) >= 2 {
			port, err = strconv.Atoi(params[1])
			if err != nil {
				log.Error(err)
				return
			}
		}

		if len(params) >= 3 {
			password = params[2]
		}

		if len(params) >= 4 {
			defaultDB, err = strconv.Atoi(params[3])
			if err != nil {
				log.Error(err)
				return
			}
		}

		if len(params) >= 5 {
			poolSize, err = strconv.Atoi(params[4])
			if err != nil {
				log.Error(err)
				return
			}
		}

		gRedis = NewRedis(host, port, password, defaultDB, poolSize)
	})

	if len(params) == 0 {
		if gRedis == nil {
			return nil, errors.New("DB is not initialized")
		}
		return gRedis, nil
	}

	return gRedis, nil
}

func NewRedis(host string, port int, password string, defaultDB int, poolSize int) *Redis {
	addr := fmt.Sprintf("%s:%d", host, port)

	rd := new(Redis)
	rd.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       defaultDB,
		PoolSize: poolSize,
	})

	return rd
}

func (o *Redis) RedisClient() *redis.Client {
	return o.client
}

func (o *Redis) Close() error {
	return o.client.Close()
}

/*
	Distributed Lock - Redlock
*/

func (o *Redis) GetLock(lockkey string, opts *LockOptions) (lock *Lock, err error) {
	lock, err = ObtainLock(o.client, lockkey, opts)
	return
}

func (o *Redis) AutoLock(lockkey string, opts *LockOptions) (func() error, error) {
	lock, err := o.GetLock(lockkey, opts)
	if err != nil {
		return nil, err
	}

	if lock == nil {
		return nil, errors.New("Lock is nil")
	}

	if ok, err := lock.Lock(); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("Obtain lock is failed")
	}

	return lock.Unlock, nil
}

/*
	HASH SET
*/

func (o *Redis) HGet(key string, field string) (string, error) {
	cmd := o.client.HGet(ctx, key, field)
	return cmd.Result()
}

func (o *Redis) HMGet(key string, fields ...string) ([]interface{}, error) {
	cmd := o.client.HMGet(ctx, key, fields...)
	return cmd.Result()
}

func (o *Redis) HGetAll(key string) (map[string]string, error) {
	cmd := o.client.HGetAll(ctx, key)
	return cmd.Result()
}

func (o *Redis) HSet(key string, field string, value string) error {
	cmd := o.client.HSet(ctx, key, field, value)
	_, err := cmd.Result()
	if err != nil {
		return err
	}

	return nil
}

func (o *Redis) HMSet(key string, fields map[string]string) error {
	cmd := o.client.HMSet(ctx, key, fields)
	_, err := cmd.Result()
	if err != nil {
		return err
	}

	return nil
}

func (o *Redis) HSetNX(key string, field string, value string) (bool, error) {
	cmd := o.client.HSetNX(ctx, key, field, value)
	ret, err := cmd.Result()
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func (o *Redis) HDel(key string, field string) error {
	cmd := o.client.HDel(ctx, key, field)
	return cmd.Err()
}

func (o *Redis) HIncrBy(key, field string, incr int64) (int64, error) {
	cmd := o.client.HIncrBy(ctx, key, field, incr)
	return cmd.Result()
}

func (o *Redis) HScan(key string, cursor uint64, match string, count int64) ([]string, error) {
	cmd := o.client.HScan(ctx, key, cursor, match, count)
	keys, _, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (o *Redis) Rename(key string, newKey string) error {
	cmd := o.client.Rename(ctx, key, newKey)
	return cmd.Err()
}

func (o *Redis) HDelAll(key string) error {
	values, err := o.HGetAll(key)

	if err != nil {
		return err
	}

	for field, _ := range values {
		if err := o.HDel(key, field); err != nil {
			return err
		}
	}

	return nil
}

func (o *Redis) HExists(key string, field string) (bool, error) {
	cmd := o.client.HExists(ctx, key, field)
	return cmd.Result()
}

func (o *Redis) HKeys(key string) ([]string, error) {
	cmd := o.client.HKeys(ctx, key)
	return cmd.Result()
}

/*
	SET
*/

func (o *Redis) Set(key string, value string) (string, error) {
	cmd := o.client.Set(ctx, key, value, 0)
	return cmd.Result()
}

func (o *Redis) SetAndExpire(key string, value string, expiration time.Duration) (string, error) {
	cmd := o.client.Set(ctx, key, value, expiration)
	return cmd.Result()
}

func (o *Redis) SetNX(key string, value string, expiration time.Duration) (bool, error) {
	cmd := o.client.SetNX(ctx, key, value, expiration)
	return cmd.Result()
}

func (o *Redis) Del(keys ...string) (int64, error) {
	cmd := o.client.Del(ctx, keys...)
	return cmd.Result()
}

func (o *Redis) Get(key string) (string, error) {
	cmd := o.client.Get(ctx, key)
	return cmd.Result()
}

func (o *Redis) MGet(keys ...string) ([]interface{}, error) {
	cmd := o.client.MGet(ctx, keys...)
	return cmd.Result()
}

// Scan Set을 조회한다.
func (o *Redis) Scan(cursor uint64, match string, count int64) ([]string, error) {
	cmd := o.client.Scan(ctx, cursor, match, count)
	keys, _, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// FullScan Set을 조회한다.
func (o *Redis) FullScan(cursor uint64, match string, count int64) ([]string, error) {
	var err error
	var result []string
	csor := cursor
	for {
		var keys []string
		cmd := o.client.Scan(ctx, csor, match, count)
		keys, csor, err = cmd.Result()
		if err != nil {
			return nil, err
		}
		result = append(result, keys...)

		if cursor == csor {
			break
		}
	}
	return result, nil
}

func (o *Redis) Keys(key string) ([]string, error) {
	cmd := o.client.Keys(ctx, key)
	return cmd.Result()
}

func (o *Redis) IncrBy(key string, value int64) (int64, error) {
	cmd := o.client.IncrBy(ctx, key, value)
	return cmd.Result()
}

func (o *Redis) Incr(key string) (int64, error) {
	cmd := o.client.Incr(ctx, key)
	return cmd.Result()
}

func (o *Redis) Decr(key string) (int64, error) {
	cmd := o.client.Decr(ctx, key)
	return cmd.Result()
}

/*
	Sorted Set
*/

// ZAdd SortedSet에 추가한다.
func (o *Redis) ZAdd(key string, members ...redis.Z) error {
	cmd := o.client.ZAdd(ctx, key, members...)
	_, err := cmd.Result()
	if err != nil {
		return err
	}
	return nil
}

// // ZAdd SortedSet에 추가한다.
// func (o *Redis) ZAddCH(key string, members ...redis.Z) error {
// 	cmd := o.client.ZAddCh(ctx, key, members...)
// 	_, err := cmd.Result()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// ZCount 범위의 개수를 Count 한다.
func (o *Redis) ZCount(key, min, max string) (int64, error) {
	cmd := o.client.ZCount(ctx, key, min, max)
	return cmd.Result()
}

// ZScore 특정 member의 스코어를 반환한다.
func (o *Redis) ZScore(key, member string) (float64, error) {
	cmd := o.client.ZScore(ctx, key, member)
	return cmd.Result()
}

// ZRank 특정 member의 랭킹을 반환한다.
func (o *Redis) ZRank(key, member string) (int64, error) {
	cmd := o.client.ZRank(ctx, key, member)
	return cmd.Result()
}

// ZRevRank 특정 member의 내림차순 랭킹을 반환한다.
func (o *Redis) ZRevRank(key, member string) (int64, error) {
	cmd := o.client.ZRevRank(ctx, key, member)
	return cmd.Result()
}

// ZRange 특정 범위의 랭킹 리스트를 반환한다.
func (o *Redis) ZRange(key string, start, stop int64) ([]string, error) {
	cmd := o.client.ZRange(ctx, key, start, stop)
	return cmd.Result()
}

// ZRangeWithScores 특정 범위의 랭킹 리스트를 반환한다.
func (o *Redis) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	cmd := o.client.ZRangeWithScores(ctx, key, start, stop)
	return cmd.Result()
}

// ZRangeByScore 특정 범위의 랭킹 리스트를 반환한다.
func (o *Redis) ZRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	cmd := o.client.ZRangeByScore(ctx, key, &opt)
	return cmd.Result()
}

// ZRangeByScore 특정 범위의 랭킹 리스트를 반환한다.
func (o *Redis) ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error) {
	cmd := o.client.ZRangeByScoreWithScores(ctx, key, &opt)
	return cmd.Result()
}

// ZRevRange 특정 범위의 내림차순 랭킹 리스트를 반환한다.
func (o *Redis) ZRevRange(key string, start, stop int64) ([]string, error) {
	cmd := o.client.ZRevRange(ctx, key, start, stop)
	return cmd.Result()
}

// ZRevRangeWithScores 특정 범위의 내림차순 랭킹 리스트를 반환한다.
func (o *Redis) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	cmd := o.client.ZRevRangeWithScores(ctx, key, start, stop)
	return cmd.Result()
}

// ZRevRangeByScore 특정 범위의 랭킹 리스트를 반환한다.
func (o *Redis) ZRevRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	cmd := o.client.ZRevRangeByScore(ctx, key, &opt)
	return cmd.Result()
}

// ZRevRangeByScoreWithScores 특정 범위의 랭킹 리스트를 반환한다.
func (o *Redis) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error) {
	cmd := o.client.ZRevRangeByScoreWithScores(ctx, key, &opt)
	return cmd.Result()
}

// // ZIncr 특정 member의 스코어를 1 더한다.
// func (o *Redis) ZIncr(key string, member redis.Z) (float64, error) {
// 	cmd := o.client.ZIncr(key, member)
// 	return cmd.Result()
// }

// ZIncrBy 특정 member의 스코어 increment를 설정한다.
func (o *Redis) ZIncrBy(key string, increment float64, member string) (float64, error) {
	cmd := o.client.ZIncrBy(ctx, key, increment, member)
	return cmd.Result()
}

// ZCard key의 개수를 반환한다.
func (o *Redis) ZCard(key string) (int64, error) {
	cmd := o.client.ZCard(ctx, key)
	return cmd.Result()
}

// ZRem members 정보를 삭제한다.
func (o *Redis) ZRem(key string, members ...interface{}) (int64, error) {
	cmd := o.client.ZRem(ctx, key, members...)
	return cmd.Result()
}

// ZremRangeByScore Score 를 기반으로 데이터를 삭제한다
func (o *Redis) ZRemRangeByScore(key string, min, max string) (int64, error) {
	cmd := o.client.ZRemRangeByScore(ctx, key, min, max)
	return cmd.Result()
}

/*
	STACK
*/

func (o *Redis) Pop(key string) (string, error) {
	cmd := o.client.LPop(ctx, key)
	return cmd.Result()
}

func (o *Redis) Push(key string, value ...interface{}) (int64, error) {
	cmd := o.client.LPush(ctx, key, value...)
	return cmd.Result()
}

/*
	QUEUE
*/
func (o *Redis) BDequeue(timeout time.Duration, key string) ([]string, error) {
	cmd := o.client.BLPop(ctx, timeout, key)
	return cmd.Result()
}

func (o *Redis) Dequeue(key string) (string, error) {
	cmd := o.client.LPop(ctx, key)
	return cmd.Result()
}

func (o *Redis) Enqueue(key string, value ...interface{}) (int64, error) {
	cmd := o.client.RPush(ctx, key, value...)
	return cmd.Result()
}

/*
	STACK and QUEUE
*/

func (o *Redis) LRange(key string, start, stop int64) ([]string, error) {
	return o.client.LRange(ctx, key, start, stop).Result()
}

func (o *Redis) LTrim(key string, start, stop int64) (string, error) {
	return o.client.LTrim(ctx, key, start, stop).Result()
}

func (o *Redis) Peek(key string) (string, error) {
	cmd := o.client.LRange(ctx, key, 0, 0)
	if len(cmd.Val()) == 1 {
		return cmd.Val()[0], nil
	}
	return "", cmd.Err()
}

/*
	PubSub
*/

func (o *Redis) Subscribe(receiveCh chan PubSubMessage, channels ...string) (chan PubSubMessage, error) {
	var err error
	if o.pubsub != nil {
		if err = o.pubsub.Subscribe(ctx, channels...); err != nil {
			return nil, err
		}
	} else {
		o.pubsub = o.client.Subscribe(ctx, channels...)
	}

	return o.startReceiveMessage(receiveCh)
}

func (o *Redis) SubscribeExpired(receiveCh chan PubSubMessage, db int) (chan PubSubMessage, error) {
	if err := o.SetNotifyExpired(); err != nil {
		return receiveCh, err
	}

	channel := fmt.Sprintf("__keyevent@%d__:expired", db)
	return o.Subscribe(receiveCh, channel)
}

func (o *Redis) startReceiveMessage(receiveCh chan PubSubMessage) (chan PubSubMessage, error) {
	if o.receiveCh != nil {
		return o.receiveCh, nil
	}

	o.receiveCh = receiveCh

	go func() {
		defer func() {
			o.receiveCh = nil
		}()

		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
			}
		}()

		for {
			if o.pubsub == nil {
				return
			}

			msg, err := o.pubsub.ReceiveMessage(ctx)
			if err != nil {
				if err.Error() != "[redis: client is closed]" {
					log.Error(err)
				}
				return
			}

			if msg == nil {
				continue
			}

			receiveCh <- msg
		}
	}()

	return o.receiveCh, nil
}

func (o *Redis) ClosePubSub() (err error) {
	defer func() {
		o.pubsub = nil
	}()

	if o.pubsub == nil {
		return
	}

	if err = o.pubsub.Close(); err != nil {
		return
	}

	return
}

func (o *Redis) Unsubscribe(channels ...string) (err error) {
	if o.pubsub == nil {
		return
	}

	if err = o.pubsub.Unsubscribe(ctx, channels...); err != nil {
		return
	}

	return
}

func (o *Redis) PSubscribe(receiveCh chan PubSubMessage, patterns ...string) (chan PubSubMessage, error) {
	var err error

	if err = o.SetNotifyKeyspaceEvent(); err != nil {
		return nil, err
	}

	if o.pubsub != nil {
		if err = o.pubsub.PSubscribe(ctx, patterns...); err != nil {
			return nil, err
		}
	} else {
		o.pubsub = o.client.PSubscribe(ctx, patterns...)
	}

	return o.startReceiveMessage(receiveCh)
}

func (o *Redis) PUnsubscribe(patterns ...string) (err error) {
	if o.pubsub == nil {
		return
	}

	if err = o.pubsub.PUnsubscribe(ctx, patterns...); err != nil {
		return
	}

	return
}

func (o *Redis) ConfigSet(parameter string, value string) error {
	cmd := o.client.ConfigSet(ctx, parameter, value)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (o *Redis) SetNotifyKeyspaceEvent() error {
	return o.ConfigSet("notify-keyspace-events", "KEAx")
}

func (o *Redis) SetNotifyExpired() error {
	return o.ConfigSet("notify-keyspace-events", "KEAx")
}

func (o *Redis) Publish(channel string, msg string) (err error) {
	cmd := o.client.Publish(ctx, channel, msg)
	if cmd.Err() != nil {
		err = cmd.Err()
		return
	}

	return nil
}

/*
	Transactions & Pipeline
*/

// Pipeline 여러 명령을 수행하기 위한 Pipeline 객체를 반환한다.
func (o *Redis) Pipeline() *redis.Pipeline {
	return o.client.Pipeline().(*redis.Pipeline)
}

// // TxPipeline Transaction이 보장되는 Pipeline 반환
// func (o *Redis) TxPipeline() *redis.Pipeline {
// 	pipeline := redis.Pipeline{}
// 	pipeline.pipe = o.client.TxPipeline().Pipeline()
// 	return &pipeline
// }

/*
	Expire & TTL
*/
func (o *Redis) Expire(key string, expiration time.Duration) error {
	cmd := o.client.Expire(ctx, key, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (o *Redis) ExpireAt(key string, tm time.Time) error {
	cmd := o.client.ExpireAt(ctx, key, tm)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (o *Redis) TTL(key string) (time.Duration, error) {
	cmd := o.client.TTL(ctx, key)

	return cmd.Result()
}

/*
	Utilities
*/

func (o *Redis) Ping() (bool, error) {
	pong, err := o.client.Ping(ctx).Result()
	if err != nil {
		return false, err
	}

	if pong == "PONG" {
		return true, nil
	}

	return false, nil
}

func IsNil(err error) bool {
	if err == redis.Nil {
		return true
	}

	return false
}

/*
	SETS
*/

// SAdd Sets에 추가한다.
func (o *Redis) SAdd(key string, members ...interface{}) error {
	cmd := o.client.SAdd(ctx, key, members...)
	_, err := cmd.Result()
	if err != nil {
		return err
	}
	return nil
}

// Sets에 속한 member 갯수를 조회한다.
func (o *Redis) SCard(key string) (int64, error) {
	cmd := o.client.SCard(ctx, key)
	return cmd.Result()
}

// 각 집합에 대한 차집합(첫번째 집합에서 두번째 집합의 Member 제거)을 구한다.
func (o *Redis) SDiff(keys ...string) ([]string, error) {
	cmd := o.client.SDiff(ctx, keys...)
	return cmd.Result()
}

// 각 집합에 대 차집합(첫번째 집합에서 두번째 집합의 Member 제거)을 구해서 새로운 집합에 저장한다.
func (o *Redis) SDiffStore(destination string, keys ...string) (int64, error) {
	cmd := o.client.SDiffStore(ctx, destination, keys...)
	return cmd.Result()
}

// 각 집합에 대 교집합을 구한다.
func (o *Redis) SInter(keys ...string) ([]string, error) {
	cmd := o.client.SInter(ctx, keys...)
	return cmd.Result()
}

// 한 교집합을 구해서 새로운 집합에 저장한다.
func (o *Redis) SInterStore(destination string, keys ...string) (int64, error) {
	cmd := o.client.SInterStore(ctx, destination, keys...)
	return cmd.Result()
}

// 집합에 member가 존재하는지 확인한다.
func (o *Redis) SIsMember(key string, member interface{}) (bool, error) {
	cmd := o.client.SIsMember(ctx, key, member)
	return cmd.Result()
}

// 집합에 데이터를 조회한다.
func (o *Redis) SMembers(key string) ([]string, error) {
	cmd := o.client.SMembers(ctx, key)
	return cmd.Result()
}

// 소스 집합의 member를 목적 집합으로 이동한다.
func (o *Redis) SMove(source, destination string, member interface{}) (bool, error) {
	cmd := o.client.SMove(ctx, source, destination, member)
	return cmd.Result()
}

// 집합에서 무작위로 member를 가져온다.
func (o *Redis) SPop(key string) (string, error) {
	cmd := o.client.SPop(ctx, key)
	return cmd.Result()
}

// 집합에서 무작위로 member N개를 가져온다.
func (o *Redis) SPopN(key string, count int64) ([]string, error) {
	cmd := o.client.SPopN(ctx, key, count)
	return cmd.Result()
}

// 집합에서 무작위로 member를 조회한다.
func (o *Redis) SRandMember(key string) (string, error) {
	cmd := o.client.SRandMember(ctx, key)
	return cmd.Result()
}

// 집합에서 무작위로 member N개를 조회한다.
func (o *Redis) SRandMemberN(key string, count int64) ([]string, error) {
	cmd := o.client.SRandMemberN(ctx, key, count)
	return cmd.Result()
}

// 집합에서 member를 삭제한다.
func (o *Redis) SRem(key string, members ...interface{}) (int64, error) {
	cmd := o.client.SRem(ctx, key, members...)
	return cmd.Result()
}

// 각 집합에 대한 합집합을 구한다.
func (o *Redis) SUnion(keys ...string) ([]string, error) {
	cmd := o.client.SUnion(ctx, keys...)
	return cmd.Result()
}

// 각 집합에 대한 합집합을 구해서 새로운 집합에 저장한다.
func (o *Redis) SUnionStore(destination string, keys ...string) (int64, error) {
	cmd := o.client.SUnionStore(ctx, destination, keys...)
	return cmd.Result()
}
