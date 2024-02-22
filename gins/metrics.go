package gins

import (
	"net/http"
	"strconv"

	"github.com/iooikaak/frame/conf/env"
	"github.com/iooikaak/frame/metadata"
	"github.com/iooikaak/frame/stat/metric/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var defaultMetricPath = "/metrics"

func (gs *Server) AddMetric(m *metrics.Metrics) {
	gs.Middleware.Use(handlerFunc(m.RegInstance()))
	gs.Router.GET(defaultMetricPath, gs.prometheusHandler(m))
}

func (gs *Server) prometheusHandler(m *metrics.Metrics) HandlerFunc {
	h := promhttp.InstrumentMetricHandler(m.RegInstance(),
		promhttp.HandlerFor(m.Gather(), promhttp.HandlerOpts{}))
	return func(c *Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func handlerFunc(reg prometheus.Registerer) HandlerFunc {
	httpReqCnt := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Name:      "request_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method", "handler", "host", "url", "from"},
	)

	ReqBusinessCode := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Name:      "business_code_request_total",
			Help:      "The number of processed HTTP requests according to business status codes and HTTP methods.",
		},
		[]string{"code", "method", "handler", "host", "url", "from"},
	)

	httpReqDur := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "http",
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
		},
		[]string{"code", "method", "handler", "host", "url", "from"},
	)

	httpResSz := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "http",
			Name:      "response_size_bytes",
			Help:      "The HTTP response sizes in bytes.",
		},
		[]string{"code", "method", "handler", "host", "url", "from"},
	)

	httpReqSz := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "http",
			Name:      "request_size_bytes",
			Help:      "The HTTP request sizes in bytes.",
		},
		[]string{"code", "method", "handler", "host", "url", "from"},
	)

	httpReqGg := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "http",
		Name:      "request_consuming_total",
		Help:      "How many HTTP requests consuming time total.",
	},
		[]string{"code", "method", "handler", "host", "url", "from"},
	)

	reg.MustRegister(httpReqCnt, httpReqDur, httpReqSz, httpResSz, ReqBusinessCode, httpReqGg)

	reqCntURLLabelMappingFn := func(c *Context) string {
		url := c.Request.URL.EscapedPath()
		if url == "" {
			url = "/"
		}

		return url
	}

	return func(c *Context) {
		if c.Request.URL.String() == defaultMetricPath {
			c.Next()
			return
		}

		url := reqCntURLLabelMappingFn(c)

		var (
			status      string
			method      = c.Request.Method
			handlerName = c.HandlerName()
			host        = c.Request.Host
			from        = c.Request.Header.Get(metadata.HttpFrom)
		)
		if from == "" {
			from = env.Hostname
		}

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			httpReqDur.WithLabelValues(status, method, handlerName, host, url, from).Observe(v)
		}))
		defer timer.ObserveDuration()

		c.Next()

		status = strconv.Itoa(c.Writer.Status())

		httpReqCnt.WithLabelValues(status, method, handlerName, host, url, from).Inc()
		httpReqSz.WithLabelValues(status, method, handlerName, host, url, from).
			Observe(float64(computeApproximateRequestSize(c.Request)))
		httpResSz.WithLabelValues(status, method, handlerName, host, url, from).
			Observe(float64(c.Writer.Size()))
		// business code
		ReqBusinessCode.WithLabelValues(strconv.Itoa(c.API.GetCode()), method, handlerName, host, url, from).Inc()
		httpReqGg.WithLabelValues(status, method, handlerName, host, url, from).Inc()
	}
}

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
