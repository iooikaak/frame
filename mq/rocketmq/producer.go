package rocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/rlog"

	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/iooikaak/frame/mq/rocketmq/config"

	"github.com/apache/rocketmq-client-go/v2"

	"github.com/iooikaak/frame/xlog"
)

type Producer struct {
	pc rocketmq.Producer
}

func NewProducer(conf *config.RocketmqConfig) (*Producer, error) {
	options := NewProducerConfig(conf)
	options = append(options, producer.WithQueueSelector(producer.NewHashQueueSelector()))
	p, err := rocketmq.NewProducer(options...)
	rlog.SetLogger(Logger())
	if err != nil {
		return nil, err
	}
	err = p.Start()
	if err != nil {
		return nil, err
	}

	return &Producer{pc: p}, nil
}

func (p *Producer) Send(ctx context.Context, topic string, body []byte, delayLevel int, tag, shardingKey string) (*primitive.SendResult, error) {
	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}
	if delayLevel > 0 {
		msg.WithDelayTimeLevel(delayLevel)
	}
	msg.WithTag(tag)
	msg.WithShardingKey(shardingKey)
	return p.pc.SendSync(ctx, msg)
}

func (p *Producer) SendAsync(ctx context.Context, topic string, body []byte, delayLevel int) error {
	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}
	if delayLevel > 0 {
		msg.WithDelayTimeLevel(delayLevel)
	}
	return p.pc.SendAsync(ctx, func(ctx context.Context, result *primitive.SendResult, e error) {
		if e != nil {
			xlog.Info("RocketMQ SendAsyncWithTopic failed result: %#v e: %v", result, e)
		}
	}, msg)
}

func (p *Producer) Stop() error {
	return p.pc.Shutdown()
}
