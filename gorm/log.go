package gorm

import "github.com/iooikaak/frame/xlog"

type LoggerFunc func(...interface{})

func (f LoggerFunc) Print(args ...interface{}) { f(args...) }

func (db *Engine) wrapLog() {
	if xlog.Logger() == nil {
		return
	}

	db.gorm.SetLogger(LoggerFunc(func(arg ...interface{}) {
		if len(arg) == 0 {
			xlog.Info(arg...)
			return
		}
		if level, ok := arg[0].(string); ok {
			switch level {
			case "sql":
				xlog.Info(arg...)
			case "warning":
				arg = arg[1:]
				xlog.Warn(arg...)
			case "error":
				arg = arg[1:]
				xlog.Error(arg...)
			}
		}

	}))
}
