package xorm

import (
	"fmt"
	"strings"
	"time"
)

var (
	defaultDatabase     = "mysql"
	ConnStrTmpl         = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s&timeout=%s&readTimeout=%s&writeTimeout=%s"
	DefaultMaxOpenConns = 200
	DefaultMaxIdleConns = 60
	DefaultMaxLeftTime  = 300 * time.Second
	DefaultTimeout      = 5 * time.Second
	DefaultWriteTimeout = 5 * time.Second
	DefaultReadTimeout  = 10 * time.Second
	TimeZone            = "Local"
	Charset             = "utf8mb4"
	MPort               = 3306
)

func New(conf *Config) (x *Engine, err error) {

	err = authConfig(conf)
	if err != nil {
		return
	}

	connStr := fmt.Sprintf(
		ConnStrTmpl,
		conf.User,
		conf.Password,
		conf.Server,
		conf.Port,
		conf.Database,
		conf.Charset,
		conf.TimeZone,
		conf.Timeout,
		conf.ReadTimeout,
		conf.WriteTimeout,
	)

	x, err = NewEngine(conf.Type, connStr)
	if err != nil {
		return
	}

	x.SetMaxIdleConns(conf.MaxIdleConns)
	x.SetMaxOpenConns(conf.MaxOpenConns)
	x.SetConnMaxLifetime(conf.MaxLeftTime)

	//registerEngineGroup(conf.Database, x)
	//SetDefaultdatabase(conf.Database)
	return
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

	if strings.TrimSpace(conf.Charset) == "" {
		conf.Charset = Charset
	}

	if strings.TrimSpace(conf.TimeZone) == "" {
		conf.TimeZone = TimeZone
	}

	return
}
