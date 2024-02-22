package config

import "time"

type Producer struct {
	SendMsgTimeout time.Duration `json:"send_msg_timeout" yaml:"send_msg_timeout"` // 发送消息超时时间
	RetryTimes     int           `json:"retry_times" yaml:"retry_times"`           // 重试次数
}
