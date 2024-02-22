package icelog

import (
	"strings"

	"github.com/iooikaak/frame/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var hostname = util.GetHostname()

// 定义日志级别
const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// 定义日志级别
var (
	levleFlags        = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	levelFlagsReverse = map[string]zapcore.Level{
		"DEBUG": zap.DebugLevel,
		"INFO":  zap.InfoLevel,
		"WARN":  zap.WarnLevel,
		"ERROR": zap.ErrorLevel,
		"FATAL": zap.FatalLevel,
	}
)

// default
var (
	defaultLogger *Logger
)

// Logger logger

type Logger struct {
	Zap      *zap.Logger
	ZapSugar *zap.Logger
	cfg      zap.Config
}

// NewLogger new logger
func NewLogger() *Logger {

	var cfg zap.Config
	cfg = zap.NewProductionConfig()
	cfg.DisableStacktrace = true
	cfg.Level.SetLevel(zap.DebugLevel)
	//cfg.EncoderConfig.EncodeTime = zapcore.CustomTimeEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//如果要写文件需要在cfg这边设置OutputPaths

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	defaultLogger = new(Logger)
	defaultLogger.Zap = logger
	defaultLogger.cfg = cfg
	defaultLogger.ZapSugar = logger.WithOptions(zap.AddCallerSkip(1)) //多包一层
	return defaultLogger
}

// SetLogPaths Set multi log path
func SetLogPaths(paths []string) {
	defaultLogger.cfg.OutputPaths = paths
}

// SetLevel SetLevel
func SetLevel(level string) {
	defaultLogger.cfg.Level.SetLevel(levelFlagsReverse[strings.ToUpper(level)])
}

// Debug global debug
func Debug(args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Debug(args...)
}

// Warn defalut wawrn
func Warn(args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Warn(args...)
}

// Info default info
func Info(args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Info(args...)
}

// Error default error
func Error(args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Error(args...)
}

// Fatal default fatal
func Fatal(args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Fatal(args...)
}

// Debugf global debug
func Debugf(fmt string, args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Debugf(fmt, args...)
}

// Warnf defalut wawrn
func Warnf(fmt string, args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Warnf(fmt, args...)
}

// Infof default info
func Infof(fmt string, args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Infof(fmt, args...)
}

// Errorf default error
func Errorf(fmt string, args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Errorf(fmt, args...)
}

// Fatalf default fatal
func Fatalf(fmt string, args ...interface{}) {
	defaultLogger.ZapSugar.Sugar().Fatalf(fmt, args...)
}

// Default default log
func Default() *Logger {
	return defaultLogger
}

func Close() {
	defaultLogger.ZapSugar.Sync()
	defaultLogger.Zap.Sync()
}
func init() {
	defaultLogger = NewLogger()
}
