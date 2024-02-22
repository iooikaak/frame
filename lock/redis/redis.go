package redisLock

import (
	"fmt"
	"time"

	"context"

	"github.com/iooikaak/frame/cache/redis/v8"
)

func New(configStr *redis.Config) (r *RedisLock, err error) {

	inst := redis.New(configStr)

	r = &RedisLock{inst}
	return
}

type RedisLock struct {
	redis *redis.Client
}

func (r *RedisLock) Lock(ctx context.Context, keySuffix string, expire int) (ok bool, err error) {
	if keySuffix == "" || expire == 0 {
		return false, fmt.Errorf("keySuffix or expire is empty")
	}

	//reply, err = redis.String(r.redis.Do(ctx, "SET", keySuffix, "", "NX"))
	ok, err = r.redis.SetNX(ctx, keySuffix, "", time.Duration(expire)*time.Second).Result()

	if err != nil {
		return false, err
	}

	return ok, err
}

func (r *RedisLock) UnLock(ctx context.Context, keySuffix string) (ok bool, err error) {
	if keySuffix == "" {
		return false, fmt.Errorf("keySuffix  is empty")
	}

	//reply, err := redis.Int(r.redis.Do(ctx, "DEL", keySuffix))
	reply, err := r.redis.Del(ctx, keySuffix).Result()
	if err != nil {
		fmt.Println("unlock发生错误", reply, err)
		return false, err
	}

	return reply == 1, err
}
