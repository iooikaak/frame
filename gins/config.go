package gins

import (
	"fmt"
	"time"
)

// Config 服务器配置
type Config struct {
	Name             string        `yaml:"name" json:"name"`                         // 服务名称，必填
	Version          string        `yaml:"version" json:"version"`                   // 服务版本，必填
	Host             string        `yaml:"host" json:"host"`                         // 域名主机
	IP               string        `yaml:"ip" json:"ip"`                             // 运行地址，必填
	BroadcastIP      string        `yaml:"broadcastIP" json:"broadcastIP"`           // 广播的运行地址，默认为：IP
	Port             int           `yaml:"port" json:"port"`                         // 运行端口，必填
	BroadcastPort    int           `yaml:"broadcastPort" json:"broadcastPort"`       // 广播的运行端口，默认为：Port
	Timeout          int           `yaml:"timeout" json:"timeout"`                   // 优雅退出时的超时机制
	Debug            bool          `yaml:"debug" json:"debug"`                       // 是否开启调试
	Pprof            bool          `yaml:"pprof" json:"pprof"`                       // 是否监控性能
	ReadTimeout      time.Duration `yaml:"readTimeout" json:"readTimeout"`           // 读超时
	WriteTimeout     time.Duration `yaml:"writeTimeout" json:"writeTimeout"`         // 写超时
	DisableAccessLog bool          `yaml:"disableAccessLog" json:"disableAccessLog"` // disable_access_log

	CenterAddr                     []string `yaml:"centerAddr" json:"centerAddr"`                                         // 注册中心地址
	CenterName                     string   `yaml:"centerName" json:"centerName"`                                         // 注册中心名称，eureka把自己也注册到了实例列表里面，用于处理集群健康检查
	RenewalIntervalInSecs          int      `yaml:"renewalIntervalInSecs" json:"renewalIntervalInSecs"`                   // 心跳间隔，单位s，默认30s
	RegistryFetchIntervalSeconds   int      `yaml:"registryFetchIntervalSeconds" json:"registryFetchIntervalSeconds"`     // 获取服务列表间隔，单位s，默认15s
	RollDiscoveriesIntervalSeconds int      `yaml:"rollDiscoveriesIntervalSeconds" json:"rollDiscoveriesIntervalSeconds"` // 滚动发现地址，单位s，默认60s
	DurationInSecs                 int      `yaml:"durationInSecs" json:"durationInSecs"`                                 // 服务过期间隔，单位s, 默认90s
	StatusUrl                      string   `yaml:"statusUrl" json:"statusUrl"`                                           // status url
	HeathUrl                       string   `yaml:"heathUrl" json:"heathUrl"`                                             // 健康检查url
}

// Addr 运行地址
func (conf *Config) Addr() string {
	return fmt.Sprintf("%s:%d", conf.IP, conf.Port)
}

// BroadcastAddr 广播的运行地址
func (conf *Config) BroadcastAddr() string {
	return fmt.Sprintf("%s:%d", conf.BroadcastIP, conf.BroadcastPort)
}
