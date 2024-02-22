package nsq

import (
	"time"

	log "github.com/iooikaak/frame/log"

	"github.com/nsqio/go-nsq"
)

// Producer nsq producer
type Producer struct {
	counter   uint64
	producers map[string]*nsq.Producer
	pushAddrs []string
	length    int
}

// NewProducer new producer
func NewProducer(addrs []string) *Producer {
	p := new(Producer)
	p.length = len(addrs)
	p.pushAddrs = addrs
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = p.length
	p.producers = make(map[string]*nsq.Producer, p.length)
	for _, v := range addrs {
		pd, err := nsq.NewProducer(v, cfg)
		if err != nil {
			log.Fatalf("connect nsq %s - %s", v, err.Error())
			return nil
		}
		p.producers[v] = pd
	}
	return p
}

// SendMsg send msg to nsq
func (p *Producer) SendMsg(topic string, body []byte) error {
	idx := p.counter % uint64(p.length)
	addr := p.pushAddrs[idx]
	p.counter++
	return p.producers[addr].Publish(topic, body)
}

// MultiSendMsg multiply send msg to nsq
func (p *Producer) SendMsgs(topic string, body [][]byte) error {
	p.counter++
	return p.producers[p.pushAddrs[p.counter%uint64(p.length)]].MultiPublish(topic, body)
}

// SendDeferredMsg 发送延时消息
func (p *Producer) SendDeferredMsg(topic string, delay time.Duration, body []byte) error {
	idx := p.counter % uint64(p.length)
	addr := p.pushAddrs[idx]
	p.counter++
	return p.producers[addr].DeferredPublish(topic, delay, body)
}

func (p *Producer) Close() error {
	for _, v := range p.producers {
		v.Stop()
	}
	return nil
}
