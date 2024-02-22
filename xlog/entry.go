package xlog

import (
	"context"
	"fmt"
	"runtime"

	"github.com/iooikaak/frame/config/env"
	"github.com/iooikaak/frame/metadata"

	"github.com/sirupsen/logrus"
)

type Entry struct {
	*logrus.Entry
	fields map[string]interface{}
}

func NewEntry(e *logrus.Entry) *Entry {
	locFields := make(map[string]interface{})

	return &Entry{
		Entry:  e,
		fields: locFields,
	}
}

func (e *Entry) WithContext(ctx context.Context) *Entry {
	if traceId := ctx.Value(metadata.HttpTraceId); traceId != nil {
		e.fields["trace_id"] = traceId
	}
	return e
}

func (e *Entry) WithEvent(event string) *Entry {
	e.fields["event"] = event
	return e
}

func (e *Entry) WithField(key string, value interface{}) *Entry {
	e.fields[key] = value
	return e
}

func (e *Entry) WithFields(fields logrus.Fields) *Entry {
	for key, value := range fields {
		e.fields[key] = value
	}
	return e
}

func (e *Entry) WithError(err error) *Entry {
	if err == nil {
		return e
	}

	if env.DeployEnv != env.DeployEnvProd {
		buf := make([]byte, 64<<10)
		buf = buf[:runtime.Stack(buf, false)]
		e.fields["stackError"] = fmt.Errorf("%s\n%s", err, buf)
	}

	e.fields["err"] = err.Error()
	return e
}

func (e *Entry) Logf(level logrus.Level, format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Logf(level, format, args...)
}

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Debugf(format, args...)
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Infof(format, args...)
}

func (e *Entry) Printf(format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Printf(format, args...)
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Warnf(format, args...)
}

func (e *Entry) Warningf(format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Warningf(format, args...)
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Errorf(format, args...)
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Fatalf(format, args...)
}

func (e *Entry) Panicf(format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Panicf(format, args...)
}

func (e *Entry) Log(level logrus.Level, args ...interface{}) {
	e.Entry.WithFields(e.fields).Log(level, args...)
}

func (e *Entry) Debug(args ...interface{}) {
	e.Entry.WithFields(e.fields).Debug(args...)
}

func (e *Entry) Info(args ...interface{}) {
	e.Entry.WithFields(e.fields).Info(args...)
}

func (e *Entry) Print(args ...interface{}) {
	e.Entry.WithFields(e.fields).Print(args...)
}

func (e *Entry) Warn(args ...interface{}) {
	e.Entry.WithFields(e.fields).Warn(args...)
}

func (e *Entry) Warning(args ...interface{}) {
	e.Entry.WithFields(e.fields).Warning(args...)
}

func (e *Entry) Error(args ...interface{}) {
	e.Entry.WithFields(e.fields).Error(args...)
}

func (e *Entry) Fatal(args ...interface{}) {
	e.Entry.WithFields(e.fields).Fatal(args...)
}

func (e *Entry) Panic(args ...interface{}) {
	e.Entry.WithFields(e.fields).Panic(args...)
}

func (e *Entry) Logln(level logrus.Level, args ...interface{}) {
	e.Entry.WithFields(e.fields).Logln(level, args...)
}

func (e *Entry) Debugln(args ...interface{}) {
	e.Entry.WithFields(e.fields).Debugln(args...)
}

func (e *Entry) Infoln(args ...interface{}) {
	e.Entry.WithFields(e.fields).Infoln(args...)
}

func (e *Entry) Println(args ...interface{}) {
	e.Entry.WithFields(e.fields).Println(args...)
}

func (e *Entry) Warnln(args ...interface{}) {
	e.Entry.WithFields(e.fields).Warnln(args...)
}

func (e *Entry) Warningln(args ...interface{}) {
	e.Entry.WithFields(e.fields).Warningln(args...)
}

func (e *Entry) Errorln(args ...interface{}) {
	e.Entry.WithFields(e.fields).Errorln(args...)
}

func (e *Entry) Fatalln(args ...interface{}) {
	e.Entry.WithFields(e.fields).Fatalln(args...)
}

func (e *Entry) Panicln(args ...interface{}) {
	e.Entry.WithFields(e.fields).Panicln(args...)
}

func (e *Entry) Tracef(format string, args ...interface{}) {
	e.Entry.WithFields(e.fields).Tracef(format, args...)
}

func (e *Entry) Trace(args ...interface{}) {
	e.Entry.WithFields(e.fields).Trace(args...)
}

func (e *Entry) Traceln(args ...interface{}) {
	e.Entry.WithFields(e.fields).Traceln(args...)
}
