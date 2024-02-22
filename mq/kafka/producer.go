package kafka

import (
	"github.com/Shopify/sarama"
	"time"
)

type KafkaProducer struct {
	sync_producer  sarama.SyncProducer
	async_producer sarama.AsyncProducer
}

func NewKafkaProducer(servers []string) (*KafkaProducer, error) {
	sync_cfg := sarama.NewConfig()
	sync_cfg.Version = sarama.V0_10_0_0
	sync_cfg.Producer.RequiredAcks = sarama.WaitForAll
	sync_cfg.Producer.Retry.Max = 10
	sync_cfg.Producer.Return.Successes = true
	if err := sync_cfg.Validate(); err != nil {
		return nil, err
	}
	sync_producer, err := sarama.NewSyncProducer(servers, sync_cfg)
	if err != nil {
		return nil, err
	}
	async_cfg := sarama.NewConfig()
	async_cfg.Version = sarama.V0_10_0_0
	async_cfg.Producer.RequiredAcks = sarama.WaitForLocal
	async_cfg.Producer.Compression = sarama.CompressionSnappy
	async_cfg.Producer.Flush.Frequency = 500 * time.Millisecond
	async_producer, err := sarama.NewAsyncProducer(servers, async_cfg)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{
		sync_producer:  sync_producer,
		async_producer: async_producer,
	}, nil
}

func (r *KafkaProducer) EncodeContent(content string) interface{} {
	return sarama.StringEncoder(content)
}

func (r *KafkaProducer) SendMessage(topic, key string, content sarama.Encoder) error {
	_, _, err := r.sync_producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     content,
		Timestamp: time.Now(),
	})
	return err
}

func (r *KafkaProducer) AsyncSendMessage(topic, key string, content sarama.Encoder) {
	r.async_producer.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     content,
		Timestamp: time.Now(),
	}
}
