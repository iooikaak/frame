# lock

## lock基础库,基于redis实现 [hms-components/lock](github.com/iooikaak/frame/lock)

> 对应底层操作：set k v ex NX

> 暴露方法： Lock() , UnLock()

# 快速开始

### 使用的例子：

```
package main

import (
	"github.com/iooikaak/frame/lock"
	"fmt"
	"context"
	"time"

	"github.com/iooikaak/frame/cache/redis"
	"github.com/iooikaak/frame/container/pool"
	xtime "github.com/iooikaak/frame/time"
)

func main() {

	redisConf := &redis.Config{
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


	ctx := context.Background()

	lockInst, err := lock.NewLock(redisConf) //实例化redis-lock
	if err != nil {
		fmt.Println("初始化失败", err)
		return
	}

	res, err := lockInst.Lock(ctx, "myLock", 2)
	if err != nil || res != true {
		fmt.Println("加锁失败", res, err)
		return
	}
	fmt.Println("加锁：", res, err)

	res2, err := lockInst.UnLock(ctx, "myLock")
	if err != nil || res2 != true {
		fmt.Println("解锁失败", res2, err)
		return
	}
	fmt.Println("解锁：", res2, err)
}
}

```

* 1.初始化实例 ``lock.NewLock(&redisConf)``

* 2.加锁操作 ``lockInst.Lock``

* 3.解锁操作 ``lockInst.UnLock``
