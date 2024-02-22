package config

const (
	// 消费模式
	ConsumerModelBroadCasting = 1 // 广播模式
	ConsumerModelClustering   = 2 // 集群模式, 默认模式

	// 消费开始处
	ConsumeFromLastOffset  = 0 // 最大offset
	ConsumeFromFirstOffset = 1 // 最小offset
	ConsumeFromTimestamp   = 2 // 从某个时间点开始消费
)

type Consumer struct {
	ConsumerModel              int    `json:"consumer_model" yaml:"consumer_model"`                                 // 消费模式
	FromWhere                  int    `json:"from_where" yaml:"from_where"`                                         // 消费开始处
	ConsumeOrderly             bool   `json:"consume_orderly" yaml:"consume_orderly"`                               // 是否是顺序消费 false 不是顺序消费 true是顺序消费
	ConsumeMessageBatchMaxSize int    `json:"consume_message_batch_max_size" yaml:"consume_message_batch_max_size"` // 批量消费数量
	RetryTimes                 int    `json:"retry_times" yaml:"retry_times"`                                       // 重试次数
	MaxReconsumeTimes          int    `json:"max_reconsume_times" yaml:"max_reconsume_times"`                       // 重复消费次数
	AutoCommit                 bool   `json:"auto_commit" yaml:"auto_commit"`                                       // 是否自动提交
	TagExpression              string `json:"tag_expression" yaml:"tag_expression"`                                 // tag表达式如： TagA || TagB
}
