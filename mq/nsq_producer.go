package mq

import (
	"log"
	"time"

	"github.com/iooikaak/frame/app"

	"github.com/nsqio/go-nsq"
)

type NsqProducer struct {
	counter    uint64
	producers  map[string]*nsq.Producer
	NsqWriters app.StringArray
}

func (p *NsqProducer) Init() {
	cfg := nsq.NewConfig()
	p.producers = make(map[string]*nsq.Producer)
	for _, addr := range p.NsqWriters {
		producer, err := nsq.NewProducer(addr, cfg)
		if err != nil {
			log.Fatalf("Connect to nsq host err")
		}
		p.producers[addr] = producer
	}
}

func (p *NsqProducer) SendNsqMsg(topic string, body []byte) error {
	idx := p.counter % uint64(len(p.NsqWriters))
	producer := p.producers[p.NsqWriters[idx]]
	err := producer.Publish(topic, body)
	p.counter++
	return err
}

func (p *NsqProducer) MultiSendNsqMsg(topic string, body [][]byte) error {
	idx := p.counter % uint64(len(p.NsqWriters))
	producer := p.producers[p.NsqWriters[idx]]
	err := producer.MultiPublish(topic, body)
	p.counter++
	return err
}

func (p *NsqProducer) SendNsqDeferredMsg(topic string, delay time.Duration, body []byte) error {
	idx := p.counter % uint64(len(p.NsqWriters))
	producer := p.producers[p.NsqWriters[idx]]
	err := producer.DeferredPublish(topic, delay, body)
	p.counter++
	return err
}
