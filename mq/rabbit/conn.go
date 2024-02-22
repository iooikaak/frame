package rabbit

import (
	"context"
	"time"

	"github.com/iooikaak/frame/xlog"

	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/streadway/amqp"
)

const delay = 50 // reconnect after delay seconds

// Connection amqp.Connection wrapper
type Connection struct {
	*amqp.Connection
	cfg      *Config
	chanPool *pool.ObjectPool
}

// Channel wrap amqp.Connection.Channel, get a auto reconnect channel
func (c *Connection) Channel() (*Channel, error) {
	channel, err := c.chanPool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	return channel.(*Channel), nil
}

// Channel wrap amqp.Connection.Channel, get a auto reconnect channel
func (c *Connection) Close() error {
	c.chanPool.Close(context.Background())
	if err := c.Connection.Close(); err != nil {
		return err
	}
	return nil
}

// Dial wrap amqp.Dial, dial and get a reconnect connection
func Dial(url string, cfg *Config) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Connection: conn,
		cfg:        cfg,
	}

	pooledConfig := pool.NewDefaultPoolConfig()
	pooledConfig.BlockWhenExhausted = cfg.Wait
	pooledConfig.MaxTotal = cfg.ChanActive

	if cfg.ChanIdleActive > 0 {
		pooledConfig.MaxIdle = cfg.ChanIdleActive
	}

	connection.chanPool = pool.NewObjectPoolWithDefaultConfig(context.Background(), &chanFactory{conn: connection, cfg: cfg})
	connection.chanPool.PreparePool(context.Background())

	go func() {
		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				xlog.WithField("param", "rabbitmq").Info("connection closed")
				break
			}
			xlog.WithField("param", "rabbitmq").Warnf("connection closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(delay * time.Millisecond)

				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Connection = conn
					xlog.WithField("param", "rabbitmq").Info("reconnect success")
					break
				}

				xlog.WithField("param", "rabbitmq").Warnf("reconnect failed, err: %v", err)
			}
		}
	}()

	return connection, nil
}
