package rocketmq

import (
	"github.com/iooikaak/frame/mq/rocketmq/config"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func NewProducerConfig(conf *config.RocketmqConfig) (options []producer.Option) {
	if len(conf.DSN) > 0 {
		options = append(options, producer.WithNsResolver(primitive.NewPassthroughResolver(conf.DSN)))
	}
	if len(conf.Namespace) > 0 {
		options = append(options, producer.WithNamespace(conf.Namespace))
	}
	if conf.Producer.SendMsgTimeout > 0 {
		options = append(options, producer.WithSendMsgTimeout(conf.Producer.SendMsgTimeout))
	}
	if conf.Producer.RetryTimes > 0 {
		options = append(options, producer.WithRetry(conf.Producer.RetryTimes))
	}
	return
}

func NewConsumerConfig(conf *config.RocketmqConfig) (options []consumer.Option) {
	if len(conf.DSN) > 0 {
		options = append(options, consumer.WithNsResolver(primitive.NewPassthroughResolver(conf.DSN)))
	}
	if len(conf.Namespace) > 0 {
		options = append(options, consumer.WithNamespace(conf.Namespace))
	}
	if conf.Consumer.ConsumerModel == config.ConsumerModelBroadCasting {
		options = append(options, consumer.WithConsumerModel(consumer.BroadCasting))
	}
	if conf.Consumer.ConsumerModel == config.ConsumerModelClustering {
		options = append(options, consumer.WithConsumerModel(consumer.Clustering))
	}
	options = append(options, consumer.WithConsumeFromWhere(consumer.ConsumeFromWhere(conf.Consumer.FromWhere)))
	options = append(options, consumer.WithConsumerOrder(conf.Consumer.ConsumeOrderly))
	if conf.Consumer.ConsumeMessageBatchMaxSize > 0 {
		options = append(options, consumer.WithConsumeMessageBatchMaxSize(conf.Consumer.ConsumeMessageBatchMaxSize))
	}
	if conf.Consumer.RetryTimes > 0 {
		options = append(options, consumer.WithRetry(conf.Consumer.RetryTimes))
	}
	if conf.Consumer.MaxReconsumeTimes > 0 {
		options = append(options, consumer.WithMaxReconsumeTimes(int32(conf.Consumer.MaxReconsumeTimes)))
	}
	options = append(options, consumer.WithAutoCommit(conf.Consumer.AutoCommit))

	return
}
