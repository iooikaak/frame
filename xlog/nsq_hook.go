package xlog

import (
	"github.com/sirupsen/logrus"
)

type AppendServiceHook struct {
	Service string
}

type AppendServiceHookConfig struct {
	Service string
}

func NewAppendServiceHook(config *AppendServiceHookConfig) *AppendServiceHook {
	return &AppendServiceHook{
		Service: config.Service,
	}
}

func (hook *AppendServiceHook) Fire(entry *logrus.Entry) error {
	entry.Data["service"] = hook.Service
	entry.Data["timestamp"] = entry.Time.Unix()
	return nil
}

func (hook *AppendServiceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
