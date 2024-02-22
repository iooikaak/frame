package redisLock

import (
	"os"
	"testing"
	"time"

	"github.com/iooikaak/frame/cache/redis/v8"
	"github.com/iooikaak/frame/container/pool"
	xtime "github.com/iooikaak/frame/time"
)

var (
	redisConf *redis.Config
)

func TestMain(m *testing.M) {
	redisConf = &redis.Config{
		Name:         "test",
		Proto:        "tcp",
		Addr:         "127.0.0.1:6379",
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}

	redisConf.Config = &pool.Config{
		Active:      20,
		Idle:        2,
		IdleTimeout: xtime.Duration(90 * time.Second),
	}
	ret := m.Run()
	os.Exit(ret)
}
