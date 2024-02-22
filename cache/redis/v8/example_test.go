package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/iooikaak/frame/container/pool"
	xtime "github.com/iooikaak/frame/time"

	"github.com/go-redis/redis/v8"
)

var r *redis.Client

func InitClient() {
	r = New(&Config{
		Config: &pool.Config{
			Active:      500,
			Idle:        800,
			WaitTimeout: xtime.Duration(800 * time.Millisecond),
			IdleTimeout: xtime.Duration(800 * time.Millisecond),
		},
		Name:         "TestName",
		Addr:         "127.0.0.1:6379",
		Proto:        "tcp",
		Auth:         "",
		DialTimeout:  xtime.Duration(800 * time.Millisecond),
		ReadTimeout:  xtime.Duration(800 * time.Millisecond),
		WriteTimeout: xtime.Duration(800 * time.Millisecond),
		DB:           0,
		SlowLog:      xtime.Duration(200 * time.Millisecond),
	})
}

//nolint
func ExampleGetString() {
	InitClient()
	a, b := r.Get(context.Background(), "test1").Result()
	fmt.Printf("--%v--%v--", a, b)
}

//nolint
func ExampleSetString() {
	InitClient()
	a, b := r.Set(context.Background(), "test1", "test1", 0).Result()
	fmt.Printf("--%v--%v--", a, b)
}

//nolint
func ExampleHashSet() {
	InitClient()
	a, b := r.HGet(context.Background(), "test_hash", "test1").Result()
	fmt.Printf("--%v--%v--", a, b)
}
