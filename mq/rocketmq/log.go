package rocketmq

import "github.com/iooikaak/frame/xlog"

type Log struct {
	xlog.ILog
}

func (log *Log) Warning(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Warn(msg)
}

func (log *Log) Level(level string) {
}

func (log *Log) Logger() *Log {
	return log
}

func (log *Log) OutputPath(path string) (err error) {
	return nil
}

func (log *Log) Debug(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Debugf(msg)
}

func (log *Log) Error(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Errorf(msg)
}

func (log *Log) Info(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Infof(msg)
}

func (log *Log) Fatal(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Fatalf(msg)
}

func Logger() *Log {
	return &Log{
		ILog: xlog.Logger(),
	}
}
