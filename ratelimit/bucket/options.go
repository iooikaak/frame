package ratelimit

import (
	"context"
	"fmt"

	"github.com/iooikaak/frame/cache/redis/v8"
)

type Options interface {
	//从集群中申请
	Apply(ctx context.Context, applyNum, tick int64) (ok bool, err error)
}

func NewDefaultOptions(opt ...RedisOption) Options {
	RedisOption := new(RedisOptions)

	for _, o := range opt {
		o(RedisOption)
	}

	o := Options(RedisOption)
	return o
}

type RedisOption func(*RedisOptions)

type RedisOptions struct {
	//每次填充的数量
	capacity int64
	//redis配置
	*redis.Client
	//要操作的key前缀
	prefix string
}

func SetRedisOptionsPrefix(key string) RedisOption {
	return func(options *RedisOptions) {
		options.prefix = key
	}
}

func SetRedisInstance(pool *redis.Client) RedisOption {
	return func(options *RedisOptions) {
		options.Client = pool
	}
}

func SetRedisOptionsCapacity(capacity int64) RedisOption {
	return func(options *RedisOptions) {
		options.capacity = capacity
	}
}

func (redis *RedisOptions) Apply(ctx context.Context, applyNum, tick int64) (ok bool, err error) {
	redisKey := fmt.Sprintf("%s%d", redis.prefix, tick)

	//同上 通过原子操作保证性能
	//conn := redis.Client.Get(ctx)
	//incrNum, err := conn.Do("INCRBY", redisKey, applyNum)
	incrNum, err := redis.Client.IncrBy(ctx, redisKey, applyNum).Result()
	if err != nil {
		return false, err
	}

	//如果大于capacity说明已经没有达到上限
	if incrNum > redis.capacity {
		return false, nil
	}

	//如果是第一次设置 添加过期时间
	if incrNum == applyNum {
		//_, err := conn.Do("EXPIRE", redisKey, DefaultRedisExpireTime)
		_, err := redis.Client.Expire(ctx, redisKey, DefaultRedisExpireTime).Result()
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
