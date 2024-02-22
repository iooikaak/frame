package bootstrap

import (
	"sync/atomic"

	"github.com/iooikaak/frame/gins"
	"github.com/iooikaak/frame/micro"
	"github.com/iooikaak/frame/stat/metric/metrics"
	"github.com/iooikaak/frame/tracer"
	"github.com/iooikaak/frame/xlog"
)

// Config struct
type Config struct {
	Log         *xlog.Config
	Tracer      *tracer.TracingConfig
	GinServer   *gins.Config
	MicroServer *micro.Options
	ServiceName string `yaml:"serviceName" json:"serviceName"`
}

type App struct {
	guardrail  int32
	conf       *Config
	httpServer *gins.Server
	// gins.Server function
	engineFunc func(*gins.Server)

	// service is a simple micro server abstraction
	service *micro.Service
	// service function
	serviceFunc func(*micro.Service)

	// metrics
	metrics *metrics.Metrics

	// tracer
	tracer *tracer.Tracer

	// Before and After funcs
	beforeStart []func()
	beforeStop  []func()
	afterStart  []func()
	afterStop   []func()
}

// New app
func New(conf Config, ops ...micro.Option) *App {
	var (
		err error
		t   *tracer.Tracer
	)
	// init log
	err = xlog.Init(conf.Log)
	if err != nil {
		panic(err)
	}

	m := metrics.New(conf.ServiceName)

	if len(conf.MicroServer.ServerName) == 0 {
		conf.MicroServer.ServerName = conf.ServiceName
	}

	if len(conf.GinServer.Name) == 0 {
		conf.GinServer.Name = conf.ServiceName
	}

	srv := micro.New(conf.MicroServer, ops...)

	if !srv.Option().DisableTrace() {
		// new tracing
		t, err = tracer.New(
			conf.ServiceName,
			*conf.Tracer,
			tracer.Logger(xlog.Logger()),
			tracer.Metrics(m.Instance()))
		if err != nil {
			panic(err)
		}
	}

	ginsServer, _ := gins.New()

	return &App{
		conf:       &conf,
		httpServer: ginsServer,
		metrics:    m,
		tracer:     t,
		service:    srv,
	}
}

// Init options
func (a *App) Init(opts ...Option) *App {
	if i := atomic.AddInt32(&a.guardrail, 1); i > 1 {
		panic("The init method can only be initialized once")
	}

	for _, opt := range opts {
		opt(a)
	}

	// init http
	//a.initGin()

	//init micro
	a.initMicro()

	return a
}

// running
func (a *App) Run() error {
	if i := atomic.LoadInt32(&a.guardrail); i == 0 {
		panic("Need to execute Init Method")
	}
	// micro start
	return a.service.Run()
}
