package mns

import (
	log "github.com/iooikaak/frame/log"
)

// Handler queue consumer
// 实现此接口，如果返回error不等与nil，则会重发消息
// 直到MNS消息保留的最大时长
// 否则代表消费成功，删除消息
type Handler interface {
	HandleMsg(MessageResp) error
	Error(error)
}

type DefaultHandler struct {
}

func (h DefaultHandler) HandleMsg(message MessageResp) error {
	log.Info(string(message.MessageBody))
	return nil
}

func (h DefaultHandler) Error(err error) {
	log.Error(err.Error())
}
