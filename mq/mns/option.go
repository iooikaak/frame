package mns

import (
	"context"
)

// ClientOptionFunc option
type ClientOptionFunc func(*Client) error

// SetHanlder 设置消费Handler
func SetHanlder(handler Handler) ClientOptionFunc {
	return func(c *Client) error {
		c.handler = handler
		return nil
	}
}

// SetRetry 设置删除消息失败重试次数
func SetBasic(retryTimes, qps, burst int) ClientOptionFunc {
	return func(c *Client) error {
		c.retry = retryTimes
		c.qps = qps
		c.burst = burst
		return nil
	}
}

// SetRetry 设置删除消息失败重试次数
func SetQueueName(queueName string) ClientOptionFunc {
	return func(c *Client) error {
		c.queueName = queueName
		return nil
	}
}

// SetContext 设置Context，用于控制消费协程退出
func SetContext(ctx context.Context) ClientOptionFunc {
	return func(c *Client) error {
		c.ctx = ctx
		return nil
	}
}
