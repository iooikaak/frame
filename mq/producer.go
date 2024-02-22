package mq

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iooikaak/frame/xlog"
)

const (
	PRODUCER_STATE_STOP ProducerState = iota //停止发送消息
	PRODUCER_STATE_RUN                       //消息发送中
	PRODUCER_STATE_WAIT                      //等待发送消息
)

type ProducerState int

func (p ProducerState) String() string {
	return producerStateToString(p)
}

func producerStateToString(s ProducerState) string {
	str := ""
	switch s {
	case PRODUCER_STATE_WAIT:
		str = "等待"
	case PRODUCER_STATE_RUN:
		str = "执行"
	case PRODUCER_STATE_STOP:
		str = "停止"
	default:
		str = fmt.Sprintf("未知[%v]", s)
	}

	return str
}

type IProducer interface {
	Init(configStr string, ilog xlog.ILog) error
	Publish(msg []byte) error
	AsyncPublish(msg []byte) error
	Stop() error
	State() ProducerState

	Name() string
	RunCount() int64
	DelayAvg() int64
}

func NewProducer(adapter, configStr string, ilog xlog.ILog) (pi IProducer, err error) {
	k := strings.ToLower(adapter)

	v, ok := mqInstance.producerAdapter[k]
	if !ok {
		err = fmt.Errorf("MQ: unknown Producer adapter %q", adapter)
		return
	}

	vo := reflect.New(v)
	pi, ok = vo.Interface().(IProducer)
	if !ok {
		err = fmt.Errorf("MQ: %q is a invalid IProducer", adapter)
		return
	}

	err = pi.Init(configStr, ilog)
	if err == nil {
		mqInstance.producerList = append(mqInstance.producerList, pi)
	}

	return
}

//获得生产者列表
func GetProducerList() []IProducer {
	return mqInstance.producerList
}
