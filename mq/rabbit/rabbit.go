package rabbit

import (
	"context"
	"time"

	pool "github.com/jolestar/go-commons-pool/v2"
)

type Config struct {
	DSN            string `yaml:"dsn" json:"dsn"`                       // 资源连接
	Active         int    `yaml:"active" json:"active"`                 // 最大连接数
	ChanActive     int    `yaml:"chanActive" json:"chanActive"`         // 最大 channel 连接数
	ChanIdleActive int    `yaml:"chanIdleActive" json:"chanIdleActive"` // 最大 channel 空闲数
	Wait           bool   `yaml:"wait" json:"wait"`                     // 是否阻塞等待
	Retry          int    `yaml:"retry" json:"retry"`                   // Publish失败是否重试
}

// RabbitMQ 基于社区提供的SDK封装连接池、重连
type MQ struct {
	connPool *pool.ObjectPool
	cfg      *Config
}

func New(cfg *Config) *MQ {
	poolConfig := pool.NewDefaultPoolConfig()
	poolConfig.MaxTotal = cfg.Active
	poolConfig.MaxIdle = cfg.Active
	poolConfig.BlockWhenExhausted = cfg.Wait
	poolConfig.TimeBetweenEvictionRuns = time.Duration(-1)

	mq := &MQ{
		connPool: pool.NewObjectPool(context.Background(), &connFactory{cfg: cfg}, poolConfig),
		cfg:      cfg,
	}
	mq.connPool.PreparePool(context.Background())
	return mq
}

// Get 获取一个连接
func (r *MQ) Get(ctx context.Context) (channel *Channel, err error) {
	defer func() {
		if e := recover(); e != nil {
			channel = nil
			err = e.(error)
		}
	}()

	pooledConn, err := r.connPool.BorrowObject(ctx)
	if err != nil {
		return nil, err
	}
	conn := pooledConn.(*Connection)
	defer func() {
		_ = r.connPool.ReturnObject(ctx, conn)
	}()

	channel, err = conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

// Release 释放一个Channel
func (r *MQ) Release(ctx context.Context, channel *Channel) error {
	return channel.conn.chanPool.ReturnObject(ctx, channel)
}

// Close 清空连接池
func (r *MQ) Close(ctx context.Context) {
	r.connPool.Close(ctx)
}
