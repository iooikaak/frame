package rabbitmq

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

func TestRabbitMQConsumer(t *testing.T) {
	maxConnNum, idleConnNum := 10, 5
	rabbitmqPool := &RabbitMQ{}
	rabbitmqPool.Init("amqp://admin:123456@10.180.18.9:5672/www.juqitech.com", maxConnNum, idleConnNum)

	fn := func() {
		conn, err := rabbitmqPool.Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		defer func() {
			if err := rabbitmqPool.Release(context.Background(), conn); err != nil {
				t.Error(err)
			}
		}()
		cn, err := conn.Channel()
		defer func() { _ = cn.Close() }()
		if err != nil {
			fmt.Println(err)
		}

		// publish
		err = cn.Publish("exchange", "key", false, false, amqp.Publishing{
			Headers:         nil,
			ContentType:     "",
			ContentEncoding: "",
			DeliveryMode:    0,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Time{},
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            nil,
		})
		if err != nil {
			fmt.Println("err", err)
		}
	}
	for i := 0; i < 10000; i++ {
		fn()
	}
}
