package xlog

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/iooikaak/frame/config/build"
	"github.com/iooikaak/frame/config/env"
	"github.com/iooikaak/frame/metadata"

	"gopkg.in/natefinch/lumberjack.v2"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	_namespace = ""
	_host      = ""
)

func init() {
	_namespace = os.Getenv("NAMESPACE")
	host, _ := os.Hostname()
	_host = host
}

type Logrus struct {
	log    *logrus.Logger
	debug  bool
	caller bool // 是否显示调用方, logrus 官方实现的caller level有bug, 所以替换成自己实现的
}

func getLogCaller(skipLevel ...int) string {
	level := 4
	if len(skipLevel) > 0 && skipLevel[0] > 0 {
		level = skipLevel[0]
	}
	_, file, line, _ := runtime.Caller(level)

	location := fmt.Sprintf("%s:%d", file, line)
	return location
}

func (l *Logrus) Infof(format string, args ...interface{}) {
	l.Logf(context.Background(), InfoLevel, format, args...)
}

func (l *Logrus) Debugf(format string, args ...interface{}) {
	l.Logf(context.Background(), DebugLevel, format, args...)
}

func (l *Logrus) Warnf(format string, args ...interface{}) {
	l.Logf(context.Background(), WarnLevel, format, args...)
}

func (l *Logrus) Errorf(format string, args ...interface{}) {
	l.Logf(context.Background(), ErrorLevel, format, args...)
}

func (l *Logrus) Fatalf(format string, args ...interface{}) {
	l.Logf(context.Background(), FatalLevel, format, args...)
}

func (l *Logrus) Printf(format string, args ...interface{}) {
	l.Logf(context.Background(), InfoLevel, format, args...)
}

func (l *Logrus) Info(args ...interface{}) {
	l.Log(context.Background(), InfoLevel, args...)
}

func (l *Logrus) Debug(args ...interface{}) {
	l.Log(context.Background(), DebugLevel, args...)
}

func (l *Logrus) Warn(args ...interface{}) {
	l.Log(context.Background(), WarnLevel, args...)
}

func (l *Logrus) Error(args ...interface{}) {
	l.Log(context.Background(), ErrorLevel, args...)
}

func (l *Logrus) Fatal(args ...interface{}) {
	l.Log(context.Background(), FatalLevel, args...)
}

func (l *Logrus) Print(args ...interface{}) {
	l.Log(context.Background(), InfoLevel, args...)
}

func (l *Logrus) Infoln(args ...interface{}) {
	l.Logln(context.Background(), InfoLevel, args...)
}

func (l *Logrus) Debugln(args ...interface{}) {
	l.Logln(context.Background(), DebugLevel, args...)
}

func (l *Logrus) Warnln(args ...interface{}) {
	l.Logln(context.Background(), WarnLevel, args...)
}

func (l *Logrus) Errorln(args ...interface{}) {
	l.Logln(context.Background(), ErrorLevel, args...)
}

func (l *Logrus) Fatalln(args ...interface{}) {
	l.Logln(context.Background(), FatalLevel, args...)
}

func (l *Logrus) Println(args ...interface{}) {
	l.Logln(context.Background(), InfoLevel, args...)
}

func (l *Logrus) IsDebug() bool {
	return l.debug
}

// Log .
func (l *Logrus) Log(ctx context.Context, level logrus.Level, args ...interface{}) {
	l.log.WithFields(l.combineFields(ctx, nil)).Log(level, args...)
}

// Logf .
func (l *Logrus) Logf(ctx context.Context, level logrus.Level, format string, args ...interface{}) {
	l.log.WithFields(l.combineFields(ctx, nil)).Logf(level, format, args...)
}

// Logln
func (l *Logrus) Logln(ctx context.Context, level logrus.Level, args ...interface{}) {
	l.log.WithFields(l.combineFields(ctx, nil)).Logln(level, args...)
}

// LogWithFields
func (l *Logrus) LogWithFields(ctx context.Context, level logrus.Level, vs map[string]interface{}, args ...interface{}) {
	l.log.WithFields(l.combineFields(ctx, vs)).Log(level, args...)
}

// LogWithFieldsFmt
func (l *Logrus) LogWithFieldsFmt(ctx context.Context, level logrus.Level, vs map[string]interface{}, format string, args ...interface{}) {
	l.log.WithFields(l.combineFields(ctx, vs)).Logf(level, format, args...)
}

// WithFields
func (l *Logrus) WithFields(vs map[string]interface{}) *Entry {
	return NewEntry(l.log.WithFields(l.combineFields(context.Background(), vs)))
}

func (l *Logrus) WithField(key string, value interface{}) *Entry {
	vs := make(map[string]interface{})
	vs[key] = value
	return NewEntry(l.log.WithFields(l.combineFields(context.Background(), vs)))
}

func (l *Logrus) WithContext(ctx context.Context) *Entry {
	return NewEntry(l.log.WithFields(l.combineFields(ctx, nil)))
}

func (l *Logrus) WithEvent(name string) *Entry {
	vs := make(map[string]interface{})
	vs["event"] = name
	return NewEntry(l.log.WithFields(l.combineFields(context.Background(), vs)))
}

func (l *Logrus) WithError(err error) *Entry {
	return NewEntry(l.log.WithFields(l.combineFields(context.Background(), nil))).WithError(err)
}

func (l *Logrus) combineFields(ctx context.Context, vsc map[string]interface{}) map[string]interface{} {
	vs := make(map[string]interface{})
	vs["zone"] = env.Cloud + "_" + env.Region + "_" + env.Zone + "_" + env.Hostname
	if build.Version != "" {
		vs["version"] = build.Version
	}
	if l.caller {
		vs["caller"] = getLogCaller()
	}
	if traceId := ctx.Value(metadata.HttpTraceId); traceId != nil {
		vs["trace_id"] = traceId
	}

	if from := ctx.Value(metadata.HttpFrom); from != nil {
		vs["from"] = from
	}
	for key, val := range vsc {
		vs[key] = val
	}

	return vs
}

// NewLogrus new logrus
func NewLogrus(conf *Config) *Logrus {
	f := &LogrusFormatter{
		Formatter: &nested.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}

	log := logrus.New()
	log.SetFormatter(f)
	log.SetLevel(DebugLevel)
	if env.DeployEnv == env.DeployEnvProd {
		log.SetLevel(InfoLevel)
	}

	if conf == nil {
		log.Out = os.Stdout
		return &Logrus{log: log, caller: true}
	}

	if conf.Service == "" {
		conf.Service = "default"
	}

	if conf.Stdout {
		log.Out = os.Stdout
		return &Logrus{log: log, caller: conf.Caller}
	}

	if conf.File != nil {
		p := path.Join(conf.File.Dir, "/"+_host+"_"+conf.File.Type+".log")

		log.Out = &lumberjack.Logger{
			Filename:   p,
			MaxSize:    conf.File.MaxSize,
			MaxBackups: conf.File.MaxBackups,
			MaxAge:     conf.File.MaxAge,
		}

		hook := NewAppendServiceHook(&AppendServiceHookConfig{Service: conf.Service})
		log.AddHook(hook)

		return &Logrus{log: log, caller: conf.Caller}
	}

	if conf.NsqConfig != nil {
		conf.NsqConfig.Service = conf.Service
		nsq, err := NewNsq(conf.NsqConfig)
		if err != nil {
			// 降级成标准输出
			log.Out = os.Stdout
			_, _ = fmt.Fprintf(os.Stdout, "NewLogrus NewNsq error(%+v)", err)
		} else {
			log.Out = nsq
		}
		hook := NewAppendServiceHook(&AppendServiceHookConfig{Service: conf.Service})
		log.AddHook(hook)
	}
	return &Logrus{log: log, caller: conf.Caller}
}
