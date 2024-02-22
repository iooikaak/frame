package tracer

import (
	"github.com/iooikaak/frame/xlog"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-lib/metrics"
)

type Option func(c *Options)

type Options struct {
	metrics           metrics.Factory
	metricInstance    *jaeger.Metrics
	logger            jaeger.Logger
	maxTagValueLength int
	queueSize         int
	poolSpans         bool
	gen128Bit         bool
}

type logWrapper struct {
	l xlog.ILog
}

func (t logWrapper) Error(msg string) {
	t.l.Error(msg)
}

func (t logWrapper) Infof(msg string, args ...interface{}) {
	t.l.Infof(msg, args...)
}

func Metrics(factory metrics.Factory) Option {
	return func(c *Options) {
		c.metrics = factory
	}
}

func Logger(logger xlog.ILog) Option {
	return func(c *Options) {
		c.logger = &logWrapper{logger}
	}
}

func PoolSpans(poolSpans bool) Option {
	return func(c *Options) {
		c.poolSpans = poolSpans
	}
}

func QueueSize(queueSize int) Option {
	return func(c *Options) {
		c.queueSize = queueSize
	}
}

func Gen128Bit(gen128Bit bool) Option {
	return func(c *Options) {
		c.gen128Bit = gen128Bit
	}
}

func MaxTagValueLength(maxTagValueLength int) Option {
	return func(c *Options) {
		c.maxTagValueLength = maxTagValueLength
	}
}

func applyOptions(options ...Option) Options {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	if opts.queueSize == 0 {
		opts.queueSize = 1000
	}

	if opts.maxTagValueLength == 0 {
		opts.maxTagValueLength = 2048
	}

	if opts.metrics != nil {
		opts.metricInstance = jaeger.NewMetrics(opts.metrics, nil)
	} else {
		opts.metricInstance = jaeger.NewMetrics(metrics.NullFactory, nil)
	}

	if opts.logger == nil {
		opts.logger = jaeger.NullLogger
	}
	return opts
}
