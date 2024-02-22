package xorm

import "time"

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
