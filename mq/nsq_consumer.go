package mq

import (
	"context"
	"log"

	"github.com/iooikaak/frame/app"

	"github.com/nsqio/go-nsq"
)

type NsqConsumer struct {
	NsqWriters app.StringArray
	NsqLookups app.StringArray
}

func (consumer *NsqConsumer) Init(topic, channel string, handler nsq.Handler, exitChan chan int) {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = len(consumer.NsqWriters)
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	concurrency := cfg.MaxInFlight
	if concurrency < 10 {
		concurrency = 10
	}
	c.AddConcurrentHandlers(handler, concurrency)

	for _, addr := range consumer.NsqWriters {
		err := c.ConnectToNSQD(addr)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, addr := range consumer.NsqLookups {
		log.Printf("lookupd addr %s", addr)
		err := c.ConnectToNSQLookupd(addr)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for {
		select {
		case <-c.StopChan:
			goto Exit
		case <-exitChan:
			c.Stop()
		}
	}
Exit:
	log.Printf("nsq consumer stopped ... topic:%s, channel:%s\n", topic, channel)
}

func (consumer *NsqConsumer) Init2(ctx context.Context, topic, channel string, handler nsq.Handler) {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = len(consumer.NsqWriters)
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	concurrency := cfg.MaxInFlight
	if concurrency < 10 {
		concurrency = 10
	}
	c.AddConcurrentHandlers(handler, concurrency)

	for _, addr := range consumer.NsqWriters {
		err := c.ConnectToNSQD(addr)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	for _, addr := range consumer.NsqLookups {
		log.Printf("lookupd addr %s", addr)
		err := c.ConnectToNSQLookupd(addr)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for {
		select {
		case <-c.StopChan:
			goto Exit
		case <-ctx.Done():
			c.Stop()
			goto Exit
		}
	}
Exit:
	log.Printf("nsq consumer stopped ... topic:%s, channel:%s\n", topic, channel)
}
