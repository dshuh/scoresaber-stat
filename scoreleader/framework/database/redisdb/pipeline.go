package redisdb

// import (
// 	"fmt"
// 	"time"

// 	redis "github.com/redis/go-redis/v9"
// )

// type Pipeline struct {
// 	pipe *redis.Pipeline
// }

// func (o *Pipeline) Exec() error {
// 	_, err := o.pipe.Exec(ctx)
// 	return err
// }

// /*
// 	HASH SET
// */

// func (o *Pipeline) HGet(key string, field string) (string, error) {
// 	cmd := o.pipe.HGet(ctx, key, field)
// 	return cmd.Result()
// }

// func (o *Pipeline) HGetAll(key string) (map[string]string, error) {
// 	cmd := o.pipe.HGetAll(ctx, key)
// 	return cmd.Result()
// }

// func (o *Pipeline) HSet(key string, field string, value string) error {
// 	cmd := o.pipe.HSet(ctx, key, field, value)
// 	_, err := cmd.Result()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (o *Pipeline) HDel(key string, field string) error {
// 	cmd := o.pipe.HDel(ctx, key, field)
// 	return cmd.Err()
// }

// func (o *Pipeline) HMDel(key string, field ...string) error {
// 	cmd := o.pipe.HDel(ctx, key, field...)
// 	return cmd.Err()
// }

// func (o *Pipeline) HDelAll(key string) error {
// 	values, err := o.HGetAll(key)

// 	if err != nil {
// 		return err
// 	}

// 	for field, _ := range values {
// 		if err := o.HDel(key, field); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (o *Pipeline) HExists(key string, field string) (bool, error) {
// 	cmd := o.pipe.HExists(ctx, key, field)
// 	return cmd.Result()
// }

// func (o *Pipeline) HKeys(key string) ([]string, error) {
// 	cmd := o.pipe.HKeys(ctx, key)
// 	return cmd.Result()
// }

// /*
// 	SET
// */

// func (o *Pipeline) Set(key string, value string) (string, error) {
// 	cmd := o.pipe.Set(ctx, key, value, 0)
// 	return cmd.Result()
// }

// func (o *Pipeline) SetAndExpire(key string, value string, expiration time.Duration) (string, error) {
// 	cmd := o.pipe.Set(ctx, key, value, expiration)
// 	return cmd.Result()
// }

// func (o *Pipeline) Get(key string) (string, error) {
// 	cmd := o.pipe.Get(ctx, key)
// 	fmt.Println(cmd.Val())
// 	return cmd.Result()
// }

// func (o *Pipeline) Del(keys ...string) (int64, error) {
// 	cmd := o.pipe.Del(ctx, keys...)
// 	return cmd.Result()
// }

// func (o *Pipeline) Keys(key string) ([]string, error) {
// 	cmd := o.pipe.Keys(ctx, key)
// 	return cmd.Result()
// }

// func (o *Pipeline) Incr(key string) (int64, error) {
// 	cmd := o.pipe.Incr(ctx, key)
// 	return cmd.Result()
// }

// func (o *Pipeline) Decr(key string) (int64, error) {
// 	cmd := o.pipe.Decr(ctx, key)
// 	return cmd.Result()
// }
