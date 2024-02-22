package bootstrap

import (
	"github.com/iooikaak/frame/micro"
	"github.com/iooikaak/frame/xlog"
)

func (a *App) initMicro() {

	var microOpt []micro.Option
	microOpt = append(microOpt, micro.HeathUrl(a.conf.GinServer.HeathUrl))
	microOpt = append(microOpt, micro.StatusUrl(a.conf.GinServer.StatusUrl))

	// aggregation
	warpOpt(&microOpt, a.beforeStart, micro.BeforeStart)
	warpOpt(&microOpt, a.beforeStop, micro.BeforeStop)
	warpOpt(&microOpt, a.afterStart, micro.AfterStart)
	warpOpt(&microOpt, a.afterStop, micro.AfterStop)
	microOpt = append(
		microOpt,
		micro.Logger(xlog.Logger()),
		micro.Metrics(a.metrics),
	)

	if a.tracer != nil {
		microOpt = append(microOpt, micro.Tracer(a.tracer.Instance()))
	}

	// init micro service
	if err := a.service.Initialize(microOpt...); err != nil {
		panic(err)
	}

	// process
	if a.serviceFunc != nil {
		a.serviceFunc(a.service)
	}
}

// batch inject micro
func warpOpt(opts *[]micro.Option, hooks []func(),
	microOption func(fn func()) micro.Option) {
	if len(hooks) > 0 {
		for _, hook := range hooks {
			h := hook
			*opts = append(*opts, microOption(func() {
				h()
			}))
		}
	}
}
