package rabbitmq

import (
	"encoding/json"
	"errors"
	"strings"
)

type Config struct {
	DSN            string `json:"dsn" yaml:"dsn"` // dsn格式连接信息
	Exchange       string `json:"exchange" yaml:"exchange"`
	Queue          string `json:"queue" yaml:"queue"`
	ConnectionsNum int    `json:"connections_num" yaml:"connections_num"` // 连接数
	RoutingKey     string `json:"routing_key" yaml:"routing_key"`
}

func newConfig(adapter, configStr string) (conf *Config, err error) {
	conf = &Config{}
	err = json.Unmarshal([]byte(configStr), conf)
	if err != nil {
		return
	}

	if conf.DSN == "" {
		err = errors.New("rabbitmq: Address cannot be empty")
	}

	if conf.ConnectionsNum == 0 {
		err = errors.New("rabbitmq: ConnectionsNum cannot be zero")
	}

	if strings.EqualFold(adapter, "consumer") && conf.Queue == "" {
		err = errors.New("rabbitmq：Queue cannot be empty")
	}

	if strings.EqualFold(adapter, "producer") && conf.Exchange == "" {
		err = errors.New("rabbitmq：Exchange cannot be empty")
	}

	if strings.EqualFold(adapter, "producer") && conf.RoutingKey == "" {
		err = errors.New("rabbitmq：RoutingKey cannot be empty")
	}

	return
}
