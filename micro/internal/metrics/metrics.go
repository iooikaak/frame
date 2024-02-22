package metrics

import (
	"context"

	"github.com/iooikaak/frame/stat/metric/metrics"

	"github.com/micro/go-micro/v2/server"
	"github.com/prometheus/client_golang/prometheus"
)

func NewHandlerWrapper(m *metrics.Metrics) server.HandlerWrapper {
	var reg = prometheus.DefaultRegisterer
	if m != nil {
		reg = m.RegInstance()
	}
	opsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "micro",
			Name:      "request_total",
			Help:      "How many go-micro requests processed, partitioned by method and status",
		},
		[]string{"method", "status"},
	)

	timeCounterSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "micro",
			Name:      "upstream_latency_microseconds",
			Help:      "Service backend method request latencies in microseconds",
		},
		[]string{"method", "status"},
	)

	timeCounterHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "micro",
			Name:      "request_duration_seconds",
			Help:      "Service method request time in seconds",
		},
		[]string{"method", "status"},
	)

	reg.MustRegister(
		opsCounter,
		timeCounterSummary,
		timeCounterHistogram,
	)
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			name := req.Endpoint()

			status := "success"
			timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
				us := v * 1000000 // make microseconds
				timeCounterSummary.WithLabelValues(name, status).Observe(us)
				timeCounterHistogram.WithLabelValues(name, status).Observe(v)
			}))
			defer timer.ObserveDuration()

			err := fn(ctx, req, rsp)
			if err != nil {
				status = "fail"
			}

			opsCounter.WithLabelValues(name, status).Inc()

			return err
		}
	}
}
