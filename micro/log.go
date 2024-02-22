package micro

import (
	"os"
	"sync"

	"github.com/iooikaak/frame/xlog"
	mlog "github.com/micro/go-micro/v2/logger"
)

type LogLevel int8

//log level form micro log
const (
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel LogLevel = iota + 1
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// InfoLevel is the default logging priority.
	// General operational entries about what's going on inside the application.
	InfoLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	ErrorLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. highest level of severity.
	FatalLevel
	//Time zone
	TimeZone = "2006-01-02 15:04:05"
)

var (
	microLogLevel = map[LogLevel]mlog.Level{
		TraceLevel: mlog.TraceLevel,
		DebugLevel: mlog.DebugLevel,
		InfoLevel:  mlog.InfoLevel,
		WarnLevel:  mlog.WarnLevel,
		ErrorLevel: mlog.ErrorLevel,
		FatalLevel: mlog.FatalLevel,
	}
)

type defaultLogger struct {
	sync.RWMutex
	opts mlog.Options
	log  xlog.ILog
}

func (l *defaultLogger) Log(level mlog.Level, v ...interface{}) {
	if !l.opts.Level.Enabled(level) || l.log == nil {
		return
	}

	switch level {
	case mlog.WarnLevel:
		l.log.Warn(v...)
	case mlog.ErrorLevel:
		l.log.Error(v...)
	case mlog.FatalLevel:
		l.log.Error(v...)
	default:
		l.log.Info(v...)
	}
}

func (l *defaultLogger) Logf(level mlog.Level, format string, v ...interface{}) {
	if !l.opts.Level.Enabled(level) || l.log == nil {
		return
	}

	switch level {
	case mlog.WarnLevel:
		l.log.Warnf(format, v...)
	case mlog.ErrorLevel:
		l.log.Errorf(format, v...)
	case mlog.FatalLevel:
		l.log.Errorf(format, v...)
	default:
		l.log.Infof(format, v...)
	}
}

func (l *defaultLogger) Options() mlog.Options {
	l.RLock()
	opts := l.opts
	l.RUnlock()
	return opts
}

func (l *defaultLogger) Init(opts ...mlog.Option) error {
	for _, o := range opts {
		o(&l.opts)
	}
	return nil
}

func (l *defaultLogger) String() string {
	return "default"
}

func (l *defaultLogger) Fields(fields map[string]interface{}) mlog.Logger {
	return l
}

func NewMicroLogger(log xlog.ILog, opts ...mlog.Option) mlog.Logger {

	options := mlog.Options{
		Fields:          make(map[string]interface{}),
		Out:             os.Stderr,
		CallerSkipCount: 2,
	}

	l := &defaultLogger{opts: options, log: log}
	if err := l.Init(opts...); err != nil {
		l.Log(mlog.FatalLevel, err)
	}

	return l
}
