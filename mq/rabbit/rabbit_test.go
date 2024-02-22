package rabbit

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

func TestProducer(t *testing.T) {
	q := New(&Config{
		DSN:        "amqp://admin:123456@10.180.18.9:5672/www.juqitech.com",
		Active:     14,
		ChanActive: 20,
		Wait:       true,
		Retry:      0,
	})
	defer q.Close(context.Background())
	fmt.Println(time.Now().String())
	fn := func(i int) {
		channel, err := q.Get(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() {
			if err := q.Release(context.Background(), channel); err != nil {
				t.Error(err)
			}
		}()

		// publish
		err = channel.Publish("amq.topic", "test", false, false, amqp.Publishing{
			Body: []byte("hello"),
		})
		if err != nil {
			t.Error(err, i)
		}
	}
	var wg sync.WaitGroup
	for x := 0; x < 20; x++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			fmt.Println(time.Now().String())
			for i := 0; i < 10000; i++ {
				fn(x)
			}
		}(x)
	}
	wg.Wait()
	fmt.Println(time.Now().String(), "end")
}

func TestConsumer(t *testing.T) {
	q := New(&Config{
		DSN:        "amqp://admin:123456@10.180.18.9:5672/www.juqitech.com",
		Active:     14,
		ChanActive: 20,
		Wait:       true,
		Retry:      0,
	})
	defer q.Close(context.Background())

	var count int64 = 0

	go func() {
		for {
			if atomic.LoadInt64(&count) == 200000 {
				q.Close(context.Background())
				time.Sleep(time.Second)
				break
			}
		}
	}()

	wg := &sync.WaitGroup{}
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cn, err := q.Get(context.Background())
			if err != nil {
				panic(err)
			}
			defer cn.Close()

			messages, _ := cn.Consume("RabbitmqTestTask", "gosdk", false, false, false, false, nil)
			for msg := range messages {
				atomic.AddInt64(&count, 1)
				_ = msg.Ack(true)
			}
		}()
	}
	wg.Wait()
	t.Log("当前消息消费数量: ", count)
}
