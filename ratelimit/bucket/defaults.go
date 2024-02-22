package ratelimit

const (
	DefaultRedisExpireTime = 2 //rediskey的默认过期时间
)

var (
	DefaultOptions = NewDefaultOptions()
)
