package bootstrap

import (
	"github.com/iooikaak/frame/gins"
	"github.com/iooikaak/frame/micro"
)

type Option func(*App)

//gin process
func HTTPService(f func(r *gins.Server)) Option {
	return func(a *App) {
		a.engineFunc = f
	}
}

//micro process
func MicroService(f func(s *micro.Service)) Option {
	return func(a *App) {
		a.serviceFunc = f
	}
}

//BeforeStart
func BeforeStart(fn func()) Option {
	return func(a *App) {
		a.beforeStart = append(a.beforeStart, fn)
	}
}

//BeforeStop
func BeforeStop(fn func()) Option {
	return func(a *App) {
		a.beforeStop = append(a.beforeStop, fn)
	}
}

//AfterStart
func AfterStart(fn func()) Option {
	return func(a *App) {
		a.afterStart = append(a.afterStart, fn)
	}
}

//AfterStop
func AfterStop(fn func()) Option {
	return func(a *App) {
		a.afterStop = append(a.afterStop, fn)
	}
}
