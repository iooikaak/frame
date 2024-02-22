package mns

import "fmt"

// 定义节点类型
type NodeType string

const (
	NODE_TYPE_DEFAULT = "default"
	NODE_TYPE_BACKUP  = "backup"
)

// QueueNode 队列节点
type QueueNode struct {
	Host         string   `toml:"host" json:"host"`
	AccessID     string   `toml:"access_id" json:"access_id"`
	AccessSecret string   `toml:"access_secret" json:"access_secret"`
	Type         NodeType `toml:"type" json:"type"` // 队列类型，backup,default
}

// OK 备份节点是否配置
func (qcfg *QueueNode) OK() bool {
	if qcfg.Host != "" && qcfg.AccessID != "" && qcfg.AccessSecret != "" {
		return true
	}
	return false
}

// SetHost set host
func (qcfg *QueueNode) SetHost(host string) {
	qcfg.Host = host
}

// SetAccessID set accessid
func (qcfg *QueueNode) SetAccessID(accessid string) {
	qcfg.AccessID = accessid
}

// SetAccessSecret set accessSecret
func (qcfg *QueueNode) SetAccessSecret(accessSecret string) {
	qcfg.AccessSecret = accessSecret
}

// SetDefault set default
func (qcfg *QueueNode) SetNodeType(tp NodeType) {
	switch tp {
	case NODE_TYPE_BACKUP, NODE_TYPE_DEFAULT:
		qcfg.Type = tp
	default:
		panic("bad node type:" + tp)
	}
}

// Config config
type Config struct {
	QPS       int         `toml:"qps" json:"qps"`
	Burst     int         `toml:"burst" json:"burst"`           // 最大Token数
	QueueName string      `toml:"queue_name" json:"queue_name"` // 消费消息的队列
	Nodes     []QueueNode `toml:"nodes" json:"nodes"`           // 消费和发送消息的节点
}

// NewConfig
func NewConfig() Config {
	return Config{QPS: 100, Burst: 10, Nodes: []QueueNode{}}
}

// Valide valide
func (config *Config) Valide() error {
	if len(config.Nodes) == 0 {
		return fmt.Errorf("forget config your mns queue ?")
	}

	if config.QPS <= 0 {
		return fmt.Errorf("qps config illegal")
	}
	if config.Burst <= 0 {
		return fmt.Errorf("burst config illegal")
	}
	if config.QueueName == "" {
		return fmt.Errorf("forget config your mns queue_name ?")
	}

	for _, v := range config.Nodes {
		if v.Host == "" {
			return fmt.Errorf("forget config your mns host ?")
		}
		if v.AccessID == "" {
			return fmt.Errorf("forget config your mns access_id ?")
		}
		if v.AccessSecret == "" {
			return fmt.Errorf("forget config your mns access_secret ?")
		}
	}
	return nil
}

// SetQPS set qps
func (config *Config) SetQPS(qps, burst int) {
	config.QPS = qps
	config.Burst = burst
}

// SetNode set node
func (config *Config) SetNode(cfg QueueNode) {
	config.Nodes = append(config.Nodes, cfg)
}

// SetQueueName set queueName
func (config *Config) SetQueueName(queueName string) {
	config.QueueName = queueName
}
