package mq

import (
	"fmt"
	"reflect"
)

var (
	mqInstance = &mq{}
)

type mq struct {
	consumerAdapter map[string]reflect.Type
	producerAdapter map[string]reflect.Type
	consumerList    []IConsumer
	producerList    []IProducer
}

// 目前用 json 的配置相对灵活,后期考虑改进
const (
	NsqConsumerConfig = "{\"address\":\"%s\",\"topic\":\"%s\",\"channel\":\"%s\"}"
	NsqProducerConfig = "{\"scheme\":\"%s\",\"address\":\"%s\",\"topic\":\"%s\"}"

	RabbitmqConsumerConfig = "{\"dsn\":\"%s\",\"queue\":\"%s\",\"connections_num\":%d}"
	RabbitmqProducerConfig = "{\"dsn\":\"%s\",\"exchange\":\"%s\",\"routing_key\":\"%s\",\"connections_num\":%d}"
)

//注册消费者适配器
func ConsumerRegister(adapter string, ci IConsumer) {
	cv := reflect.ValueOf(ci)
	ct := reflect.Indirect(cv).Type()

	if mqInstance.consumerAdapter == nil {
		mqInstance.consumerAdapter = make(map[string]reflect.Type)
	}
	mqInstance.consumerAdapter[adapter] = ct
}

//注册发布者适配器
func ProducerRegister(adapter string, pi IProducer) {
	pv := reflect.ValueOf(pi)
	pt := reflect.Indirect(pv).Type()

	if mqInstance.producerAdapter == nil {
		mqInstance.producerAdapter = make(map[string]reflect.Type)
	}

	mqInstance.producerAdapter[adapter] = pt
}

func Exit() error {
	var (
		err    error
		errMsg string
	)
	for _, p := range mqInstance.producerList {
		err = p.Stop()
		if err != nil {
			if len(errMsg) > 0 {
				errMsg += ","
			}
			errMsg += fmt.Sprintf("%q", err)
		}
	}

	for _, c := range mqInstance.consumerList {
		err = c.Stop()
		if err != nil {
			if len(errMsg) > 0 {
				errMsg += ","
			}
			errMsg += fmt.Sprintf("%q", err)
		}
	}

	if len(errMsg) > 0 {
		err = fmt.Errorf("MQ: Exit error %s", errMsg)
	}

	return err
}
