package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Environment string

var (
	ENV_DEV  Environment = "dev"
	ENV_QA   Environment = "qa"
	ENV_STAG Environment = "stag"
	ENV_PROD Environment = "prod"
)

func (env Environment) String() string {
	return string(env)
}

func (env Environment) Production() bool {
	return env == ENV_PROD
}

func (env Environment) Stag() bool {
	return env == ENV_STAG
}

func (env Environment) QA() bool {
	return env == ENV_QA
}

func (env Environment) Dev() bool {
	return env == ENV_DEV
}

// TIDBCfg mysql config
type TIDBCfg struct {
	Charset              string `toml:"charset" json:"charset" yaml:"charset"`
	Database             string `toml:"database" json:"database" yaml:"database"`
	Host                 string `toml:"host" json:"host" yaml:"host"`
	MysqlConnMaxLifeTime int32  `toml:"mysqlConnMaxLifeTime" json:"mysqlConnMaxLifeTime" yaml:"mysqlConnMaxLifeTime"`
	MysqlIdle            int32  `toml:"mysqlIdle" json:"mysqlIdle" yaml:"mysqlIdle"`
	MysqlMaxConnections  int32  `toml:"mysqlMaxConnections" json:"mysqlMaxConnections" yaml:"mysqlMaxConnections"`
	Password             string `toml:"password" json:"password" yaml:"password"`
	TimeZone             string `toml:"timeZone" json:"timeZone" yaml:"timeZone"`
	User                 string `toml:"user" json:"user" yaml:"user"`
}

// MysqlCfg mysql config
type MysqlCfg struct {
	Host struct {
		Read  string `toml:"read" json:"read" yaml:"read"`
		Write string `toml:"write" json:"write" yaml:"write"`
	} `toml:"host" json:"host" yaml:"host"`

	Port    int    `toml:"port" json:"port" yaml:"port"`
	User    string `toml:"user" json:"user" yaml:"user"`
	Psw     string `toml:"password" json:"password" yaml:"password"`
	DbName  string `toml:"dbName" json:"dbName" yaml:"dbName"`
	LogMode bool   `toml:"logMode" json:"logMode" yaml:"logMode"`
}

// RedisCfg config
type RedisCfg struct {
	Addr string `toml:"address" json:"address" yaml:"address"`
	Psw  string `toml:"password" json:"password" yaml:"password"`
	DBNo int    `toml:"dbNo" json:"dbNo" yaml:"dbNo"`
}

// BaseCfg 服务基础配置
type BaseCfg struct {
	Etcd   EtcdCfg   `toml:"etcdConfig" json:"etcdConfig" yaml:"etcdConfig"`
	Zipkin ZipkinCfg `toml:"zipkinConfig" json:"zipkinConfig" yaml:"zipkinConfig"`
	Staff  StaffCfg  `toml:"staffConfig" json:"staffConfig" yaml:"staffConfig"`
}

// ZipkinCfg Zipkin配置
type ZipkinCfg struct {
	EndPoints string `json:"endPoints" yaml:"endPoints"`
}

// StaffCfg 服务监控人员
type StaffCfg struct {
	Name        string `json:"name" toml:"name" yaml:"name"`       // 服务名称
	Email       string `json:"email" toml:"email" yaml:"email"`    // 通知邮箱
	MobilePhone string `json:"mobile" toml:"mobile" yaml:"mobile"` // 通知手机号 ，暂未实现
}

// EtcdCfg 对应配置文件中关于etcd配置内容
type EtcdCfg struct {
	EndPoints []string      `json:"endpoints" toml:"endpoints" yaml:"endpoints"`
	User      string        `json:"user" toml:"user" yaml:"user"`
	Psw       string        `json:"password" toml:"password" yaml:"password"`
	Timeout   time.Duration `json:"timeout" toml:"timeout" yaml:"timeout"`
}

type NsqCfg struct {
	Topic   string   `json:"topic" toml:"topic" yaml:"topic"`
	Writers []string `json:"writers" toml:"writers" yaml:"writers"`
	Lookups []string `json:"lookups" toml:"lookups" yaml:"lookups"`
}

type EsCfg struct {
	Hosts    []string `json:"hosts" toml:"hosts" yaml:"hosts"`
	UserName string   `json:"username" toml:"username" yaml:"username"`
	Password string   `json:"password" toml:"password" yaml:"password"`
}

type ShenceCfg struct {
	Url         string `json:"url" toml:"url" yaml:"url"`
	Project     string `json:"project" toml:"project" yaml:"project"`
	Timeout     int    `json:"timeout" toml:"timeout" yaml:"timeout"`
	LoggingPath string `json:"loggingPath" toml:"loggingPath" yaml:"loggingPath"`
}

type KafkaCfg struct {
	Servers []string `json:"servers" toml:"servers" yaml:"servers"`
}

type AliyunOssCfg struct {
	EndPoint        string `json:"endpoint" yaml:"endpoint"`
	AccessKeyID     string `json:"accessKeyID" yaml:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`
}

// Parseconfig parse json config
// out must be pointer
func Parseconfig(filepath string, out interface{}) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(out)
	if err != nil {
		panic(err)
	}
}

type OssCfg struct {
	EndPoint  string `json:"endPoint" yaml:"endPoint"`
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

//IConfig 配置接口对象
type IConfig interface {
	Type() ConfigType
	Body() ([]byte, error)
}

//NewConfig 创建新的配置对象
func NewConfig(typ ConfigType, configStr string) (IConfig, error) {
	switch typ {
	case CONFIG_TYPE_CONSUL:
		return NewConsulConfig(configStr)
	case CONFIG_TYPE_FILE:
		return NewFileConfig(configStr)
	}

	return nil, errors.New("unknown ConfigType")
}
