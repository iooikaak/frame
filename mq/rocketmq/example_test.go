package rocketmq

import (
	"context"
	"fmt"
	"testing"

	"github.com/iooikaak/frame/xlog"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/iooikaak/frame/mq/rocketmq/config"
)

func TestProducer(t *testing.T) {
	ctx := context.Background()
	topic := "TOPIC-DEV_A-NM-MTC_CHANNEL"
	//topic := "INFRA-TOPIC-DEV-GOTEST"
	p, err := NewProducer(&config.RocketmqConfig{
		DSN: []string{"10.1.2.215:9876"},
		Producer: config.Producer{
			RetryTimes: 2,
		},
	})
	if err != nil {
		xlog.Errorf("NewProducer p: %#v err: %v", p, err)
		return
	}
	for i := 0; i < 100; i++ {
		shardingKey := fmt.Sprintf("test-%d", i)
		r, err := p.Send(ctx, topic, []byte(fmt.Sprintf("TEST-VALUE--%d", i)), 0, "TagA", shardingKey)
		if err != nil {
			xlog.Errorf("Producer Send r: %#v err: %v", r, err)
			return
		}
	}
}

func TestNewPushConsumer(t *testing.T) {
	a := make(chan struct{})
	//topic := "INFRA-TOPIC-DEV-GOTEST"
	topic := "TOPIC-DEV_A-NM-MTC_CHANNEL"
	group := "INFRA-TOPIC-DEV-GOTESTGROUP"
	p, err := NewPushConsumer(&config.RocketmqConfig{
		DSN:      []string{"10.1.2.215:9876"},
		Consumer: config.Consumer{},
	}, group)
	if err != nil {
		xlog.Errorf("NewPushConsumer p: %#v err: %v", p, err)
		return
	}
	selector := ConsumerMessageSelector{
		//Type:       consumer.TAG,
		//Expression: "TagA || TagC",
	}
	err = p.RegisterHandle(topic, selector, CallBack)
	if err != nil {
		xlog.Errorf("Consumer RegisterHandle err: %v", err)
		return
	}
	err = p.Start()
	if err != nil {
		xlog.Errorf("Consumer Start err: %v", err)
		return
	}
	defer func() { _ = p.Stop() }()
	<-a
}

func CallBack(msgs []*primitive.MessageExt, ctx *primitive.ConsumeConcurrentlyContext) (ConsumeResult, error) {
	for _, msg := range msgs {
		xlog.Infof("Consumer msg: %#v body: %+v queue1111: %d", msg, string(msg.Body), msg.Queue.QueueId)
	}
	return ConsumeSuccess, nil
}
