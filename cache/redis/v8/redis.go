package redis

import (
	"time"

	xtime "github.com/iooikaak/frame/time"

	"github.com/iooikaak/frame/container/pool"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	*pool.Config `yaml:"pool" json:"pool"`

	Name         string         `yaml:"name" json:"name"` // redis name, for trace
	Proto        string         `yaml:"proto" json:"proto"`
	Addr         string         `yaml:"addr" json:"addr"`
	Auth         string         `yaml:"auth" json:"auth"`
	DialTimeout  xtime.Duration `yaml:"dialTimeout" json:"dialTimeout"`
	ReadTimeout  xtime.Duration `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout xtime.Duration `yaml:"writeTimeout" json:"writeTimeout"`
	DB           int            `yaml:"db" json:"db"`
	SlowLog      xtime.Duration `yaml:"slowLog" json:"slowLog"`
}

//New 实例化新的redis v8
func New(conf *Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Auth,
		DB:           conf.DB,
		WriteTimeout: time.Duration(conf.WriteTimeout),
		ReadTimeout:  time.Duration(conf.ReadTimeout),
		IdleTimeout:  time.Duration(conf.IdleTimeout),
		MinIdleConns: conf.Idle,
		PoolSize:     conf.Active, //缩放连接数
		PoolTimeout:  time.Duration(conf.WaitTimeout),
		DialTimeout:  time.Duration(conf.DialTimeout),
	})
	rdb.PoolStats()
	rdb.AddHook(&OpenTracingHook{
		cfg:    conf,
		status: rdb.PoolStats(),
	})
	return rdb
}

func NewScript(script string) *Script {
	return redis.NewScript(script)
}
