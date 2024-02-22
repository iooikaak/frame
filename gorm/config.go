package gorm

import (
	"fmt"
	"time"
)

type Config struct {
	Alias        string        `yaml:"alias" json:"alias"`
	Type         string        `yaml:"type" json:"type"`
	Server       string        `yaml:"server" json:"server"`
	Port         int           `yaml:"port" json:"port"`
	Database     string        `yaml:"database" json:"database"`
	User         string        `yaml:"user" json:"user"`
	Password     string        `yaml:"password" json:"password"`
	MaxIdleConns int           `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns int           `yaml:"maxOpenConns" json:"maxOpenConns"`
	Charset      string        `yaml:"charset" json:"charset"`
	TimeZone     string        `yaml:"timezone" json:"timezone"`
	MaxLeftTime  time.Duration `yaml:"maxLeftTime" json:"maxLeftTime"`
	ReadTimeout  time.Duration `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
	Timeout      time.Duration `yaml:"timeout" json:"timeout"`
}

func authConfig(conf *Config) (err error) {

	if len(conf.Type) == 0 {
		conf.Type = defaultDatabase
	}

	if conf.Port == 0 {
		conf.Port = MPort
	}

	if len(conf.User) == 0 || len(conf.Password) == 0 {
		err = fmt.Errorf("User or  Password is empty")
		return
	}

	if len(conf.Server) == 0 {
		err = fmt.Errorf("server addr is empty")
		return
	}

	if len(conf.Database) == 0 {
		err = fmt.Errorf("database is empty")
		return
	}

	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = DefaultMaxIdleConns
	}

	if conf.MaxLeftTime == 0 {
		conf.MaxLeftTime = DefaultMaxLeftTime
	}

	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = DefaultMaxOpenConns
	}

	if conf.Timeout == 0 {
		conf.Timeout = DefaultTimeout
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = DefaultReadTimeout
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = DefaultWriteTimeout
	}

	return
}
