package tracer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/iooikaak/frame/xlog"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/thrift"
	j "github.com/uber/jaeger-client-go/thrift-gen/jaeger"
	"time"
)

const (
	DialTimeout      = 3 * time.Second
	ReadTimeout      = 10 * time.Second
	WriteTimeout     = 10 * time.Second
	JaegerTopic      = "jaegerTopics"
	BatchSize        = 128000
	Frequency        = 500
	BatchMaxMessages = 10
)

type KafkaSender struct {
	client    sarama.AsyncProducer
	host      []string
	batchSize int
	spans     []*j.Span
	process   *j.Process
}

func NewKafkaTransport(hostPort ...string) (jaeger.Transport, error) {

	if len(hostPort) == 0 {
		return nil, fmt.Errorf("kafka broket ip is empty")
	}

	k := &KafkaSender{
		batchSize: 20,
		spans:     []*j.Span{},
		host:      hostPort,
	}

	err := k.Init()
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (c *KafkaSender) Init() (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Flush.Bytes = BatchSize
	config.Producer.Flush.Frequency = time.Millisecond * Frequency
	config.Producer.Flush.MaxMessages = BatchMaxMessages
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Net.DialTimeout = DialTimeout
	config.Net.ReadTimeout = ReadTimeout
	config.Net.WriteTimeout = WriteTimeout

	c.client, err = sarama.NewAsyncProducer(c.host, config)
	if err != nil {
		return err
	}
	xlog.Errorf("kafka broker [%v] topic [%v]", c.host, JaegerTopic)
	go c.handleError()
	return
}

func (c *KafkaSender) Append(span *jaeger.Span) (int, error) {
	if c.process == nil {
		c.process = jaeger.BuildJaegerProcessThrift(span)
	}
	jSpan := jaeger.BuildJaegerThrift(span)
	c.spans = append(c.spans, jSpan)
	if len(c.spans) >= c.batchSize {
		return c.Flush()
	}
	return 0, nil
}

func (c *KafkaSender) Flush() (int, error) {
	count := len(c.spans)
	if count == 0 {
		return 0, nil
	}
	err := c.send(c.spans)
	c.spans = c.spans[:0]
	return count, err
}

func (c *KafkaSender) Close() error {
	return nil
}

func (c *KafkaSender) send(spans []*j.Span) error {
	batch := &j.Batch{
		Spans:   spans,
		Process: c.process,
	}

	body, err := serializeThrift(batch)
	if err != nil {
		xlog.Errorf("producer serializeThriftError %v", err)
		return err
	}

	c.client.Input() <- &sarama.ProducerMessage{
		Topic: JaegerTopic,
		Value: sarama.ByteEncoder(body.Bytes()),
	}

	return nil
}

func serializeThrift(obj thrift.TStruct) (*bytes.Buffer, error) {
	t := thrift.NewTMemoryBuffer()
	p := thrift.NewTBinaryProtocolConf(t, &thrift.TConfiguration{})
	if err := obj.Write(context.TODO(), p); err != nil {
		return nil, err
	}
	return t.Buffer, nil
}

func (c *KafkaSender) handleError() {
	var (
		err *sarama.ProducerError
		ok  bool
	)
	for {

		err, ok = <-c.client.Errors()
		if !ok {
			xlog.Errorf("producer ProducerError has be closed, break the handleError goroutine")
			return
		}

		if err != nil {
			xlog.Errorf("producer message error, partition:%d offset:%d key:%v valus:%s error(%v)\n", err.Msg.Partition, err.Msg.Offset, err.Msg.Key, err.Msg.Value, err.Err)
		}

	}
}
