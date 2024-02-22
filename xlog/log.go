package xlog

import (
	"context"
)

var (
	h ILog
)

// Config log config.
type Config struct {
	Service   string     `yaml:"service"`
	Host      string     `yaml:"host"`
	Caller    bool       `yaml:"caller"`    // 是否显式调用
	Stdout    bool       `yaml:"stdout"`    // 优先级第一
	File      *DirConfig `yaml:"file"`      // 优先级第二
	NsqConfig *NsqConfig `yaml:"nsqconfig"` // 优先级第三
}

func init() {
	cfg := &Config{
		Service: "default",
		Stdout:  true,
		Caller:  true,
	}
	h = NewLogrus(cfg)
}

func Init(conf *Config) (err error) {
	h = NewLogrus(conf)
	return nil
}

func Empty() bool {
	return h == nil
}

func Logger() ILog {
	return h
}

func Infof(format string, args ...interface{}) {
	h.Logf(context.Background(), InfoLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	h.Logf(context.Background(), DebugLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	h.Logf(context.Background(), WarnLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	h.Logf(context.Background(), ErrorLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	h.Logf(context.Background(), FatalLevel, format, args...)
}

func Infoln(args ...interface{}) {
	h.Logln(context.Background(), InfoLevel, args...)
}

func Debugln(args ...interface{}) {
	h.Logln(context.Background(), DebugLevel, args...)
}

func Warnln(args ...interface{}) {
	h.Logln(context.Background(), WarnLevel, args...)
}

func Errorln(args ...interface{}) {
	h.Logln(context.Background(), ErrorLevel, args...)
}

func Fatalln(args ...interface{}) {
	h.Logln(context.Background(), FatalLevel, args...)
}

func Debug(args ...interface{}) {
	h.Log(context.Background(), DebugLevel, args...)
}

func Info(args ...interface{}) {
	h.Log(context.Background(), InfoLevel, args...)
}

func Warn(args ...interface{}) {
	h.Log(context.Background(), WarnLevel, args...)
}

func Error(args ...interface{}) {
	h.Log(context.Background(), ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	h.Log(context.Background(), FatalLevel, args...)
}

func InfoWithFields(vs map[string]interface{}, args ...interface{}) {
	h.LogWithFields(context.Background(), InfoLevel, vs, args...)
}

func WarnWithFields(vs map[string]interface{}, args ...interface{}) {
	h.LogWithFields(context.Background(), WarnLevel, vs, args...)
}

func ErrorWithFields(vs map[string]interface{}, args ...interface{}) {
	h.LogWithFields(context.Background(), ErrorLevel, vs, args...)
}

func InfoWithFieldsFmt(vs map[string]interface{}, format string, args ...interface{}) {
	h.LogWithFieldsFmt(context.Background(), InfoLevel, vs, format, args...)
}

func WarnWithFieldsFmt(vs map[string]interface{}, format string, args ...interface{}) {
	h.LogWithFieldsFmt(context.Background(), WarnLevel, vs, format, args...)
}

func ErrorWithFieldsFmt(vs map[string]interface{}, format string, args ...interface{}) {
	h.LogWithFieldsFmt(context.Background(), ErrorLevel, vs, format, args...)
}

func WithFields(vs map[string]interface{}) *Entry {
	return h.WithFields(vs)
}

func WithField(key string, value interface{}) *Entry {
	return h.WithField(key, value)
}

func WithContext(ctx context.Context) *Entry {
	return h.WithContext(ctx)
}

func WithEvent(name string) *Entry {
	return h.WithEvent(name)
}

func WithError(err error) *Entry {
	return h.WithError(err)
}

func Warning(format string, args ...interface{}) {
	h.Logf(context.Background(), WarnLevel, format, args...)
}

//func Level(level string) {
//	h.SetLevel(level)
//}
