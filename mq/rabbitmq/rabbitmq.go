package rabbitmq

import (
	"context"
	"fmt"

	"github.com/iooikaak/frame/mq"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	pool "github.com/jolestar/go-commons-pool/v2"
)

func init() {
	mq.ConsumerRegister("rabbitmq", &Consumer{})
	mq.ProducerRegister("rabbitmq", &Producer{})
}

// RabbitPoolFactory 创建Rabbit连接池工厂类
type RabbitPoolFactory struct {
	dsn string
}

// MakeObject 创建对象方法
func (f *RabbitPoolFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	fmt.Println("make执行")
	conn, err := rabbitmq.Dial(f.dsn)
	return pool.NewPooledObject(conn), err
}

// DestroyObject 回收对象方法
func (f *RabbitPoolFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	fmt.Println("删除执行")
	v := object.Object.(*rabbitmq.Connection)
	if !v.IsClosed() {
		return v.Close()
	}
	return nil
}

// ValidateObject 验证对象有效性
func (f *RabbitPoolFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	fmt.Println("验证执行")
	v := object.Object.(*rabbitmq.Connection)
	return v.IsClosed()
}

// ActivateObject 激活对象
func (f *RabbitPoolFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// v := object.Object.(*rabbitmq.Connection)
	return nil
}

// PassivateObject 钝化对象
func (f *RabbitPoolFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

// RabbitMQ 基于社区提供的SDK封装连接池、重连
type RabbitMQ struct {
	factory *RabbitPoolFactory
	PoolCtx context.Context
	Pool    *pool.ObjectPool
}

// Init 初始化
func (r *RabbitMQ) Init(dsn string, maxNum int, maxIdleNum int) {
	r.factory = &RabbitPoolFactory{}
	r.factory.dsn = dsn

	r.PoolCtx = context.TODO()
	r.Pool = pool.NewObjectPoolWithDefaultConfig(r.PoolCtx, r.factory)
	r.Pool.Config.MaxTotal = maxNum
	r.Pool.Config.MaxIdle = maxIdleNum
}

// Get 获取一个连接
func (r *RabbitMQ) Get(ctx context.Context) (conn *rabbitmq.Connection, err error) {
	defer func() {
		if e := recover(); e != nil {
			conn = nil
			err = e.(error)
		}
	}()

	c, err := r.Pool.BorrowObject(ctx)
	if err != nil {
		return nil, err
	}
	return c.(*rabbitmq.Connection), err
}

// Release 释放一个连接
func (r *RabbitMQ) Release(ctx context.Context, conn *rabbitmq.Connection) error {
	return r.Pool.ReturnObject(ctx, conn)
}

// Close 清空连接池
func (r *RabbitMQ) Close(ctx context.Context) {
	r.Pool.Close(ctx)
}
