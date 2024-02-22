package xlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ILog interface {
	Logf(context.Context, logrus.Level, string, ...interface{})
	Logln(context.Context, logrus.Level, ...interface{})
	Log(context.Context, logrus.Level, ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Print(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Debug(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Warn(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Fatal(args ...interface{})
	IsDebug() bool
	LogWithFields(c context.Context, level logrus.Level, vs map[string]interface{}, args ...interface{})
	LogWithFieldsFmt(c context.Context, level logrus.Level, vs map[string]interface{}, format string, args ...interface{})
	WithFields(vs map[string]interface{}) *Entry
	WithField(key string, value interface{}) *Entry
	WithContext(ctx context.Context) *Entry
	WithEvent(name string) *Entry
	WithError(err error) *Entry
}

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = logrus.PanicLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = logrus.FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = logrus.ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = logrus.WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = logrus.InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = logrus.DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = logrus.TraceLevel
)
