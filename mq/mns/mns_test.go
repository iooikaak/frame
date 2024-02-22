package mns

import (
	"testing"
)

func TestBadMNS(t *testing.T) {
	q := NewQueue(Config{
		QPS:       1,
		Burst:     5,
		QueueName: "play-mns-test",
		Nodes: []QueueNode{
			QueueNode{
				Host:         "https://iooikaak.mns.cn-hangzhou.aliyuncs.com",
				AccessID:     "LTAIptdGl4O3wVq2",
				AccessSecret: "0YsvICco2gLngkRv8ScECvpeuFaASE",
				Type:         "backup",
			},
			QueueNode{
				Host:         "https://iooikaak.mns.cn-qingdao.aliyuncs.com",
				AccessID:     "LTAIptdGl4O3wVq2",
				AccessSecret: "0YsvICco2gLngkRv8ScECvpeuFaASE",
				Type:         "default",
			},
		},
	}, SetHanlder(DefaultHandler{}))

	for i := 0; i < 100000; i++ {
		if err := q.SendMsg("play-mns-test", Message{
			MessageBody:  []byte("123456"),
			DelaySeconds: 0,
			Priority:     8,
		}); err != nil {
			t.Error(err.Error())
		} else {
			t.Log("success")
		}
		// time.Sleep(time.Second * 10)
	}
}
