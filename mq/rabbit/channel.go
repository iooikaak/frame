package rabbit

import (
	"sync/atomic"
	"time"

	"github.com/iooikaak/frame/xlog"

	"github.com/streadway/amqp"
)

// Channel amqp.Channel wrapper
type Channel struct {
	*amqp.Channel
	conn   *Connection
	cfg    *Config
	closed int32
}

// IsClosed indicate closed by developer
func (ch *Channel) IsClosed() bool {
	return atomic.LoadInt32(&ch.closed) == 1
}

// Close ensure closed flag set
func (ch *Channel) Close() error {
	if ch.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&ch.closed, 1)

	return ch.Channel.Close()
}

// Consume warp amqp.Channel.Consume, the returned delivery will end only when channel closed by developer
func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := ch.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				xlog.WithField("param", "rabbitmq").Warnf("consume failed, err: %v", err)
				time.Sleep(delay * time.Millisecond)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(delay * time.Millisecond)

			if ch.IsClosed() {
				close(deliveries)
				break
			}
		}
	}()

	return deliveries, nil
}

func (ch *Channel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) (err error) {
	retry := ch.cfg.Retry + 1 // 最少执行一次
	for retry > 0 {
		err = ch.Channel.Publish(exchange, key, mandatory, immediate, msg)
		if err != nil {
			retry--
			time.Sleep(time.Millisecond * 100)
			continue
		}
		break
	}
	return err
}
