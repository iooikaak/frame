package rocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/rlog"

	"github.com/iooikaak/frame/mq/rocketmq/config"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type ConsumeResult int

const (
	ConsumeSuccess ConsumeResult = iota
	ConsumeRetryLater
	Commit
	Rollback
	SuspendCurrentQueueAMoment
)

type PushConsumer struct {
	pc rocketmq.PushConsumer
}

type PushCallback func(msgs []*MessageExt, ctx *ConsumeConcurrentlyContext) (ConsumeResult, error)

func NewPushConsumer(conf *config.RocketmqConfig, groupName string) (pc *PushConsumer, err error) {
	options := NewConsumerConfig(conf)
	options = append(options, consumer.WithGroupName(groupName))
	c, err := rocketmq.NewPushConsumer(options...)
	if err != nil {
		return nil, err
	}
	rlog.SetLogger(Logger())

	return &PushConsumer{
		pc: c,
	}, nil
}

func (pc *PushConsumer) RegisterHandle(topic string, s ConsumerMessageSelector, f PushCallback) (err error) {
	err = pc.pc.Subscribe(topic, s, func(ctx context.Context, msgs ...*primitive.MessageExt) (result consumer.ConsumeResult, e error) {
		concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
		r, e := f(msgs, concurrentCtx)
		return consumer.ConsumeResult(r), e
	})
	if err != nil {
		return err
	}
	return
}

func (pc *PushConsumer) Start() error {
	return pc.pc.Start()
}

func (pc *PushConsumer) Stop() error {
	return pc.pc.Shutdown()
}
