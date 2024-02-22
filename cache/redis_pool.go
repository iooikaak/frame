package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisPool struct {
	*redis.Pool
}

func NewRedisPool(server, password string, maxIdle int, idleTimeout time.Duration) *RedisPool {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &RedisPool{Pool: pool}
}
