package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/iooikaak/frame/mq"
	"github.com/iooikaak/frame/xlog"
)

type Consumer struct {
	connPool *RabbitMQ
	isInit   bool
	conf     *Config
	handler  mq.ConsumerHandlerFunc
	state    mq.ConsumerState
	err      error
	runCount int64 //消费统计
	delayAvg int64 //平均延迟
	delaySum int64 //累计延迟

	// 因为内部没有实现,所以需要消费时封装处理
	stopChan   chan struct{}  // 关闭指令
	exitedChan chan struct{}  // 关闭成功后结果
	wg         sync.WaitGroup // 等待组,等待所有消费者安全退出

	ilog xlog.ILog
}

// configStr:
func (c *Consumer) Init(configStr string, ilog xlog.ILog) (err error) {
	if c.isInit {
		return
	}

	conf, err := newConfig("consumer", configStr)
	if err != nil {
		err = fmt.Errorf("rabbitmq：Consumer config Error：%q", err)
		return
	}

	c.connPool = &RabbitMQ{}
	// 最大空闲连接等于最大连接的50%加1
	c.connPool.Init(conf.DSN, conf.ConnectionsNum, conf.ConnectionsNum*2+1)
	c.conf = conf
	c.isInit = true
	c.exitedChan = make(chan struct{}, 1)
	c.stopChan = make(chan struct{}, 1)
	c.ilog = ilog

	return
}
func (c *Consumer) Handler(handler mq.ConsumerHandlerFunc) error {
	c.handler = handler
	return nil
}

func (c *Consumer) Start() error {
	c.state = mq.CONSUMER_STATE_WAIT

	initError := make(chan error, c.conf.ConnectionsNum+1)
	connectedChan := make(chan struct{}, c.conf.ConnectionsNum+1)

	// 启动消费池
	c.Consume(initError, connectedChan)

	// 全部链接成功后返回,否则等待失败结果或超时
	for i := 0; i < c.conf.ConnectionsNum; i++ {
		select {
		case <-connectedChan:
			c.ilog.Infoln("consumer started!")
			continue
		case err := <-initError:
			c.state = mq.CONSUMER_STATE_STOP
			_ = c.Stop(err)
			return err
		}
	}

	c.state = mq.CONSUMER_STATE_RUN

	return nil
}

// Consume 启动消费池
func (c *Consumer) Consume(initError chan error, connectedChan chan struct{}) {
	for i := 0; i < c.conf.ConnectionsNum; i++ {
		c.wg.Add(1)
		go func(stopChan <-chan struct{}, initErrorChan <-chan error) {
			defer func() {
				if e := recover(); e != nil {
					c.ilog.Errorln("rabbitmq consumer error: %s", e.(error).Error())
					initError <- e.(error)
				}
				c.wg.Done()
			}()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			conn, err := c.connPool.Get(ctx)
			if err != nil {
				initError <- err
				return
			}

			defer func() {
				// 连接归还连接池
				_ = c.connPool.Release(context.Background(), conn)
			}()

			cn, err := conn.Channel()
			if err != nil {
				panic(err)
			}
			defer func() { _ = cn.Close() }() // 屏蔽了返回错误,因为只有当已经关闭的时候才会返回错误

			messages, err := cn.Consume(c.conf.Queue, "", false, false, false, false, nil)
			if err != nil {
				panic(err)
			}

			// 通知主进程已经链接成功开始获取消息
			connectedChan <- struct{}{}

			for {
				select {
				case <-initErrorChan:
					return
				case <-stopChan:
					return
				case msg, ok := <-messages:
					if !ok {
						c.ilog.Infoln("消费时发现消息已经关闭,正常退出")
						return
					}
					// 会返回错误,但是业务错误不处理
					_ = c.handler(c, &mq.Message{
						Body:      msg.Body,
						ID:        msg.MessageId,
						Timestamp: msg.Timestamp.Unix(),
						Attempts:  int(msg.MessageCount),
						Object:    msg,
					})
				}
			}
		}(c.stopChan, initError)
	}

	go func() {
		// 等待所有消费者安全退出,发送已经退出指令
		c.wg.Wait()
		c.exitedChan <- struct{}{}
	}()
}

func (c *Consumer) Stop(...error) error {
	if !c.isInit {
		return errors.New(`rabbitmq：Consumer Stop error "must first call the Init func"`)
	}

	close(c.stopChan)
	<-c.exitedChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	c.connPool.Close(ctx)

	c.ilog.Infoln("rabbitmq: Consumer closed")
	return nil
}

func (c *Consumer) State() mq.ConsumerState {
	return c.state
}

func (c *Consumer) Error() error {
	return c.err
}

func (c *Consumer) Name() string {
	return fmt.Sprintf("consumer-%s", c.conf.Queue)
}

func (c *Consumer) RunCount() int64 {
	return c.runCount
}
func (c *Consumer) DelayAvg() int64 {
	return c.delayAvg
}
