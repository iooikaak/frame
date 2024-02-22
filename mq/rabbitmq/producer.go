package rabbitmq

import (
	"context"
	"fmt"

	"github.com/iooikaak/frame/mq"
	"github.com/iooikaak/frame/xlog"

	"github.com/streadway/amqp"
)

// FIXME: 其实 connections 不需要太大, 只需要创建一个 channel pool 就可以了,待优化

type Producer struct {
	connPool *RabbitMQ
	isInit   bool
	conf     *Config
	state    mq.ProducerState
	err      error
	runCount int64 // 消费统计
	delayAvg int64 // 平均延迟
	delaySum int64 // 累计延迟

	ilog xlog.ILog
}

// 实例化rabbitmq
func InitProducer(conf *Config) (producer *Producer, err error) {
	producer = &Producer{}

	producer.connPool = &RabbitMQ{}
	producer.connPool.Init(conf.DSN, conf.ConnectionsNum, conf.ConnectionsNum/2+1)
	producer.conf = conf
	producer.isInit = true

	return
}

func (p *Producer) Init(configStr string, ilog xlog.ILog) error {
	if p.isInit {
		return nil
	}

	conf, err := newConfig("producer", configStr)
	if err != nil {
		return fmt.Errorf("rabbitmq：Producer config Error：%s, %s", err.Error(), configStr)
	}

	p.connPool = &RabbitMQ{}
	// 最大空闲连接等于最大连接的50%加1
	p.connPool.Init(conf.DSN, conf.ConnectionsNum, conf.ConnectionsNum*2+1)
	p.conf = conf
	p.isInit = true
	p.ilog = ilog

	return nil
}

// 异步发送消息
func (p *Producer) AsyncPublish(msg []byte) (err error) {
	err = fmt.Errorf("Does not support asynchronous mode")
	return
}

func (p *Producer) Publish(msg []byte) error {
	ctx := context.Background()
	conn, err := p.connPool.Get(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = p.connPool.Release(ctx, conn)
	}()
	cn, err := conn.Channel()
	if err != nil {
		return err
	}
	defer cn.Close()
	rabbitMsg := amqp.Publishing{
		Body: msg,
	}
	err = cn.Publish(p.conf.Exchange, p.conf.RoutingKey, false, false, rabbitMsg)
	return err
}

func (p *Producer) PublishByExchangeRoutingKey(ctx context.Context, exchange, routingKey string, msg []byte) error {
	conn, err := p.connPool.Get(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = p.connPool.Release(ctx, conn)
	}()
	cn, err := conn.Channel()
	if err != nil {
		return err
	}
	defer cn.Close()
	rabbitMsg := amqp.Publishing{
		Body: msg,
	}
	err = cn.Publish(exchange, routingKey, false, false, rabbitMsg)
	return err
}

func (p *Producer) Stop() error {
	return nil
}

func (p *Producer) State() mq.ProducerState {
	return 0
}

func (p *Producer) Name() string {
	return ""
}
func (p *Producer) RunCount() int64 {
	return 0
}
func (p *Producer) DelayAvg() int64 {
	return 0
}
