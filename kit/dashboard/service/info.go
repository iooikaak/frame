package service

import (
	"github.com/iooikaak/frame/mq"
	"github.com/iooikaak/frame/task"
)

// Info 服务信息
type Info struct {
	bill map[string]interface{}
}

// NewInfo 实例
func NewInfo() *Info {
	info := &Info{}
	info.bill = make(map[string]interface{})
	return info
}

// Bill 清单
func (info *Info) Bill() (map[string]interface{}, error) {
	return info.bill, nil
}

// ServiceName 服务名称
func (info *Info) ServiceName(name string) *Info {
	info.bill["name"] = name
	return info
}

// IsLeader 是否Leader
func (info *Info) IsLeader(isLeader bool) *Info {
	info.bill["isLeader"] = isLeader
	return info
}

// Task 任务
func (info *Info) Task(t *task.TaskManager) *Info {

	if t == nil {
		return info
	}

	info.bill["task"] = t.Info()

	return info
}

// MQ 消息队列
func (info *Info) MQ() *Info {

	mqInfo := &infoMQ{
		Version: "",
	}

	producerList := mq.GetProducerList()
	producerListLen := len(producerList)

	mqInfo.Producer.Start = producerListLen
	mqInfo.Producer.Items = make([]infoMQProducerItem, producerListLen)

	for i := 0; i < producerListLen; i++ {
		mqInfo.Producer.Items[i].Name = producerList[i].Name()
		mqInfo.Producer.Items[i].Status = producerList[i].State().String()
		mqInfo.Producer.Items[i].RunCount = producerList[i].RunCount()
		mqInfo.Producer.Items[i].DelayAvg = producerList[i].DelayAvg()
	}

	consumerList := mq.GetConsumerList()
	consumerListLen := len(consumerList)

	mqInfo.Consumer.Start = consumerListLen
	mqInfo.Consumer.Items = make([]infoMQConsumerItem, consumerListLen)

	for i := 0; i < consumerListLen; i++ {
		errMsg := "nil"
		if consumerList[i].Error() != nil {
			errMsg = consumerList[i].Error().Error()
		}

		mqInfo.Consumer.Items[i].Name = consumerList[i].Name()
		mqInfo.Consumer.Items[i].Status = consumerList[i].State().String()
		mqInfo.Consumer.Items[i].RunCount = consumerList[i].RunCount()
		mqInfo.Consumer.Items[i].DelayAvg = consumerList[i].DelayAvg()
		mqInfo.Consumer.Items[i].Error = errMsg
	}

	info.bill["mq"] = mqInfo

	return info
}

type infoMQ struct {
	Version  string
	Producer struct {
		Start int
		Items []infoMQProducerItem
	}
	Consumer struct {
		Start int
		Items []infoMQConsumerItem
	}
}

type infoMQConsumerItem struct {
	Name     string
	Status   string
	RunCount int64
	DelayAvg int64
	Error    string
}

type infoMQProducerItem struct {
	Name     string
	Status   string
	RunCount int64
	DelayAvg int64
}
