package nsq

import (
	"context"

	log "github.com/iooikaak/frame/log"

	"github.com/nsqio/go-nsq"
)

// Config nsq config
type Config struct {
	topic     string
	pushAddrs []string
	lookups   []string
	channal   string
}

// NewConfig new nsq config
func NewConfig(topic, channal string, pushAddrs, lookups []string) *Config {
	var c = new(Config)
	c.topic = topic
	c.pushAddrs = pushAddrs
	c.lookups = lookups
	c.channal = channal
	return c
}

// Consumer consumer
type Consumer struct {
	pushAddrs []string
	lookups   []string
	topic     string
	channal   string
	handler   nsq.Handler
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewConsumer new nsq consumer
func NewConsumer(cfg *Config, handler nsq.Handler) *Consumer {
	c := new(Consumer)
	c.pushAddrs = cfg.pushAddrs
	c.lookups = cfg.lookups
	c.topic = cfg.topic
	c.channal = cfg.channal
	c.handler = handler
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

// Run run consumer
func (consumer *Consumer) Run() {
	cfg := nsq.NewConfig()
	// MaxInFlight NSQ 建议大于等于nsqd节点数
	cfg.MaxInFlight = len(consumer.pushAddrs)
	c, err := nsq.NewConsumer(consumer.topic, consumer.channal, cfg)
	if err != nil {
		panic(err.Error())
	}

	concurrency := cfg.MaxInFlight
	if concurrency < 10 {
		concurrency = 10
	}

	c.AddConcurrentHandlers(consumer.handler, concurrency)

	for _, addr := range consumer.pushAddrs {
		err := c.ConnectToNSQD(addr)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	for _, addr := range consumer.lookups {
		log.Infof("lookupd addr %s", addr)
		err := c.ConnectToNSQLookupd(addr)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	for {
		select {
		case <-consumer.ctx.Done():
			c.Stop()
			log.Infof("nsq consumer stopped ... topic:%s, channel:%s\n", consumer.topic, consumer.channal)
			return
		}
	}
}

// Stop stop consumer
func (consumer *Consumer) Stop() {
	consumer.cancel()
}
