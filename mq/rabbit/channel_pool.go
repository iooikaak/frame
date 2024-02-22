package rabbit

import (
	"context"
	"time"

	"github.com/iooikaak/frame/xlog"

	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/streadway/amqp"
)

// RabbitChanFactory 创建Rabbit连接池工厂类
type chanFactory struct {
	conn *Connection
	cfg  *Config
}

// MakeObject 创建对象方法
func (f *chanFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	ch, err := f.conn.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		Channel: ch,
		conn:    f.conn,
		cfg:     f.cfg,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.IsClosed() {
				xlog.WithField("param", "rabbitmq").Warn("channel closed")
				channel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			xlog.WithField("param", "rabbitmq").Warnf("channel closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(delay * time.Millisecond)

				ch, err := f.conn.Connection.Channel()
				if err == nil {
					xlog.WithField("param", "rabbitmq").Info("channel recreate success")
					channel.Channel = ch
					break
				}

				xlog.WithField("param", "rabbitmq").Warnf("channel recreate failed, err: %v", err)
			}
		}

	}()

	return pool.NewPooledObject(channel), err
}

// DestroyObject 回收对象方法
func (f *chanFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	v := object.Object.(*Channel)
	if !v.IsClosed() {
		return v.Close()
	}
	return nil
}

// ValidateObject 验证对象有效性
func (f *chanFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return object.Object.(*Channel).IsClosed()
}

// ActivateObject 激活对象
func (f *chanFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// v := object.Object.(*Channel)
	return nil
}

// PassivateObject 钝化对象
func (f *chanFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
