package config

type RocketmqConfig struct {
	DSN       []string `json:"dsn" yaml:"dsn"`             // dsn格式连接信息
	Namespace string   `json:"namespace" yaml:"namespace"` // 命名空间配置
	Producer  Producer `json:"producer" yaml:"producer"`   // 生产者配置
	Consumer  Consumer `json:"consumer" yaml:"consumer"`   // 消费者配置
}
