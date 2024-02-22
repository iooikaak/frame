package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/iooikaak/frame/cache/redis/v8"
	ratelimit "github.com/iooikaak/frame/ratelimit/bucket"

	"gopkg.in/yaml.v3"
)

type Configs struct {
	Redis *redis.Config `yaml:"redis"`
}

func InitRedis() *redis.Client {

	redisConfig := new(Configs)

	yamlFile, err := filepath.Abs("./redis.yaml")

	if err != nil {
		fmt.Printf("InitRedis err:%v\n", err)
		return nil
	}

	yamlRead, err := ioutil.ReadFile(yamlFile)

	if err != nil {
		fmt.Printf("InitRedis err:%v\n", err)
		return nil
	}

	err = yaml.Unmarshal(yamlRead, redisConfig)
	if err != nil {
		fmt.Printf("InitRedis err:%v\n", err)
		return nil
	}

	RedisPool := redis.New(redisConfig.Redis)

	return RedisPool
}

func TestRedis(t *testing.T) {
	ctx := context.Background()
	p := InitRedis()

	reply, err := p.HGet(ctx, "Hkey", "fld").Result()
	fmt.Printf("reply:%v err:%v\n", reply, err)

	reply, err = p.Get(ctx, "test").Result()
	fmt.Printf("reply:%v err:%v\n", reply, err)
}

func NewBucket(applyNum, Capacity int64, redisPrefix, action string) *ratelimit.Bucket {
	prefixKey := redisPrefix + action
	bucketOptions := ratelimit.NewDefaultOptions(
		ratelimit.SetRedisOptionsPrefix(prefixKey),
		ratelimit.SetRedisOptionsCapacity(Capacity),
		ratelimit.SetRedisInstance(InitRedis()),
	)

	bucket := ratelimit.NewBucket(
		ratelimit.SetApplyNum(applyNum),
		ratelimit.SetAvailableTokens(applyNum),
		ratelimit.SetOptions(bucketOptions),
		ratelimit.SetCapacity(applyNum),
	)

	return bucket
}

func TestRateLimit(t *testing.T) {
	bucker := NewBucket(5, 40, "redidPrefix", "test")
	ctx := context.Background()
	t.Parallel()

	var testBucker int64

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("TestRateLimit:%d", i), func(t *testing.T) {
			ok, err := bucker.Take(ctx, 1)
			if err != nil {
				t.Errorf("TestRateLimit err:%v", err)
				return
			}
			if ok {
				atomic.AddInt64(&testBucker, 1)
			}
		})
	}

	time.Sleep(2 * time.Second)
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("TestRateLimit:%d", i), func(t *testing.T) {
			ok, err := bucker.Take(ctx, 1)
			if err != nil {
				t.Errorf("TestRateLimit err:%v", err)
				return
			}
			if ok {
				atomic.AddInt64(&testBucker, 1)
			}
		})
	}

	time.Sleep(2 * time.Second)

	if testBucker != 85 {
		t.Errorf("error occur testBucker:%v\n", testBucker)
		return
	}

	t.Logf("testBucker:%v\n", testBucker)

}

func BenchmarkRateLimit(b *testing.B) {
	bucker := NewBucket(5, 40, "redidPrefix", "test")

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := bucker.Take(ctx, 1); err != nil {
			b.Error(err)
		}
	}

}
