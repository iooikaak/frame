package lock

import (
	"context"
	"fmt"

	"github.com/iooikaak/frame/cache/redis/v8"
	redisLock "github.com/iooikaak/frame/lock/redis"
)

type ILock interface {
	Lock(ctx context.Context, keySuffix string, expire int) (bool, error)
	UnLock(ctx context.Context, keySuffix string) (bool, error)
}

//创建锁实例 默认redis锁
func NewLock(config interface{}, lockType ...string) (l ILock, err error) {
	var lockTypeFlag string

	if len(lockType) > 0 {
		lockTypeFlag = lockType[0]
	}

	switch lockTypeFlag {
	default:
		if v, ok := config.(*redis.Config); ok {
			l, err = redisLock.New(v)
			return
		}

	}

	err = fmt.Errorf("unknown lock type, allowing redis")
	return
}
