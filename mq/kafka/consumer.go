package kafka

import (
	"context"
	"sync/atomic"

	log "github.com/iooikaak/frame/log"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type Handler interface {
	HandleMessage(message *sarama.ConsumerMessage) error
}

type HandlerFunc func(message *sarama.ConsumerMessage) error

// HandleMessage implements the Handler interface
func (h HandlerFunc) HandleMessage(m *sarama.ConsumerMessage) error {
	return h(m)
}

type KafkaConsumer struct {
	consumer        *cluster.Consumer
	runningHandlers int32
	ctx             context.Context
}

func NewKafkaConsumer(groupID string, servers, topic []string, ctx context.Context) (*KafkaConsumer, error) {
	cfg := cluster.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(servers, groupID, topic, cfg)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{
		consumer: consumer,
		ctx:      ctx,
	}, nil
}

func (r *KafkaConsumer) AddHandler(handler Handler) {
	atomic.AddInt32(&r.runningHandlers, 1)
	go func() {
		for {
			select {
			case msg, ok := <-r.consumer.Messages():
				if !ok {
					goto exit
				}
				err := handler.HandleMessage(msg)
				if err != nil {
					log.Warn(err.Error())
					continue
				}
				r.consumer.MarkOffset(msg, "")
			case err, ok := <-r.consumer.Errors():
				if ok {
					log.Warn(err.Error())
				}
			case ntf, ok := <-r.consumer.Notifications():
				if ok {
					log.Infof("Kafka consumer rebalance:%+v", ntf)
				}
			case <-r.ctx.Done():
				goto exit
			}
		}
	exit:
		log.Warnf("stopping handler, running handlers %d", atomic.AddInt32(&r.runningHandlers, -1))
	}()
}
