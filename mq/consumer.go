package mq

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/iooikaak/frame/xlog"
)

const (
	CONSUMER_STATE_STOP ConsumerState = iota //停止接收消息
	CONSUMER_STATE_RUN                       //消息处理中
	CONSUMER_STATE_WAIT                      //等待消息投递
)

type ConsumerState int

func (c ConsumerState) String() string {
	return consumerStateToString(c)
}

func consumerStateToString(s ConsumerState) string {
	str := ""
	switch s {
	case CONSUMER_STATE_WAIT:
		str = "等待"
	case CONSUMER_STATE_RUN:
		str = "执行"
	case CONSUMER_STATE_STOP:
		str = "停止"
	default:
		str = fmt.Sprintf("未知[%v]", s)
	}

	return str
}

type ConsumerHandlerFunc func(IConsumer, *Message) error

type IConsumer interface {
	Init(configStr string, ilog xlog.ILog) error
	Handler(handler ConsumerHandlerFunc) error
	Start() error
	Stop(...error) error
	State() ConsumerState
	Error() error

	Name() string
	RunCount() int64
	DelayAvg() int64
}

func NewConsumer(adapter, configStr string, ilog xlog.ILog, handler ConsumerHandlerFunc) (ci IConsumer, err error) {
	if handler == nil {
		err = errors.New("MQ：handler func cannot be nil")
		return
	}

	v, ok := mqInstance.consumerAdapter[adapter]
	if !ok {
		err = fmt.Errorf("MQ: unknown Consumer adapter %q", adapter)
		return
	}

	vo := reflect.New(v)
	ci, ok = vo.Interface().(IConsumer)
	if !ok {
		err = fmt.Errorf("MQ: %q is a invalid IConsumer", adapter)
		return
	}

	err = ci.Init(configStr, ilog)
	if err != nil {
		return
	}

	err = ci.Handler(handler)
	if err != nil {
		return
	}

	mqInstance.consumerList = append(mqInstance.consumerList, ci)

	return
}

//获得消费者列表
func GetConsumerList() []IConsumer {
	return mqInstance.consumerList
}
