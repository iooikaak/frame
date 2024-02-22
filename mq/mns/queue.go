package mns

import (
	"fmt"
	"net/http"

	log "github.com/iooikaak/frame/log"
)

// Queue MNS 队列
type Queue struct {
	c   *Client
	cfg Config
	// 当前节点类型
	currentNodeType string

	// 备用节点
	backupQueueNode QueueNode
	// 默认使用的MNS队列
	defaultQueueNode QueueNode
}

// NewQueue new queue
func NewQueue(cfg Config, opts ...ClientOptionFunc) *Queue {
	q := new(Queue)
	if len(cfg.Nodes) > 2 {
		panic("backup mns node only support one")
	}
	var defaultQ bool
	for _, v := range cfg.Nodes {

		switch v.Type {
		case NODE_TYPE_DEFAULT:
			q.defaultQueueNode = v
			defaultQ = true
		case NODE_TYPE_BACKUP:
			q.backupQueueNode = v
		default:
			panic("bad node type:" + v.Type)
		}
	}

	if !defaultQ {
		panic("forget config your default mns ?")
	}

	// 默认使用default 节点
	q.currentNodeType = NODE_TYPE_DEFAULT
	opts = append(opts, SetQueueName(cfg.QueueName), SetBasic(5, cfg.QPS, cfg.Burst))
	q.c = NewClient(opts...)
	if q.c.handler != nil {
		go q.c.recvMessages(q.defaultQueueNode)
		go q.c.recvMessages(q.backupQueueNode)
	}
	return q
}

// SendMsg 发送消息
func (q *Queue) SendMsg(queueName string, message Message) error {
	resource := fmt.Sprintf("queues/%s/%s", queueName, "messages")
	if message.Priority == 0 {
		message.Priority = 8
	}

	var sResp *http.Response
	var err error
	switch q.currentNodeType {
	case NODE_TYPE_DEFAULT:
		sResp, err := q.c.Send("POST", nil, &message, resource, q.defaultQueueNode)
		if err == nil {
			sResp.Body.Close()
			return nil
		}
		// 重试 5 次失败，转移到备份节点
		q.currentNodeType = NODE_TYPE_BACKUP
		log.Info("使用备份节点:", q.backupQueueNode.Host)

	case NODE_TYPE_BACKUP:
		if !q.backupQueueNode.OK() { // 使用备份发送
			return errMnsBackupNodeNotOpen
		}
		sResp, err = q.c.Send("POST", nil, &message, resource, q.backupQueueNode)
		if err == nil {
			sResp.Body.Close()
		}
	}
	return err
}

// BatchSendMsg 批量发送消息,返回error代表不成功，error为错误信息
func (q *Queue) BatchSendMsg(queueName string, message ...Message) error {
	resource := fmt.Sprintf("queues/%s/%s", queueName, "messages")

	var batchMsg BatchMessage
	batchMsg.Messages = make([]Message, len(message))
	for i, v := range message {
		if batchMsg.Messages[i].Priority == 0 {
			batchMsg.Messages[i].Priority = 8
		}
		batchMsg.Messages[i] = v
	}

	var sResp *http.Response
	var err error
	switch q.currentNodeType {
	case NODE_TYPE_DEFAULT:
		sResp, err = q.c.Send("POST", nil, &batchMsg, resource, q.defaultQueueNode)
		if err == nil {
			sResp.Body.Close()
			return nil
		}
		q.currentNodeType = NODE_TYPE_BACKUP
		log.Info("使用备份节点:", q.backupQueueNode.Host)

	case NODE_TYPE_BACKUP:
		if !q.backupQueueNode.OK() { // 使用备份发送
			return errMnsBackupNodeNotOpen
		}
		sResp, err = q.c.Send("POST", nil, &batchMsg, resource, q.backupQueueNode)
		if err == nil {
			sResp.Body.Close()
		}
	}
	return err
}
