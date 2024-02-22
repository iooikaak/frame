package tracer

import (
	"context"

	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"
	"net/http"
)

type TraceOperator interface {
	SetTracingData(k, v string)
	GetTracingData(k string) string
	SetTag(k string, v interface{})
	Inject(r *http.Request) error
	Extract(header http.Header) (opentracing.SpanContext, error)
	GetTracingId() string
	GetSpanId() string
	GetParentId() string
	GetSpan() opentracing.Span
	Finish()
}

const (
	jaegerTracing = 1
)

var tracingType int

//hms
func StartTraceOperator(r **http.Request, SpanName string, extract ...bool) (t TraceOperator, err error) {

	if !opentracing.IsGlobalTracerRegistered() {
		return nil, fmt.Errorf("please init  trace.New")
	}

	var isExtract bool
	if len(extract) > 0 && extract[0] {
		isExtract = extract[0]
	}

	//TODO 可以增加任意实现啦opentracing协议的链路客户端
	switch tracingType {
	default:
		return startJaegerTracingOperator(r, SpanName, isExtract)
	}
}

func StartTraceCtx(r context.Context, SpanName string) (context.Context, TraceOperator, error) {

	if !opentracing.IsGlobalTracerRegistered() {
		return r, nil, fmt.Errorf("please init  trace.New")
	}

	//TODO 可以增加任意实现啦opentracing协议的链路客户端
	switch tracingType {
	default:
		return startTracingCtx(r, SpanName)
	}

}

type Tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
	conf   *TracingConfig
}

func New(serviceName string, tracingConfig TracingConfig, op ...Option) (t *Tracer, err error) {

	var (
		config = tracingConfig
	)

	if len(serviceName) == 0 {
		err = fmt.Errorf("serviceName is empty ")
		return
	}

	if len(config.SamplerType) == 0 {
		err = fmt.Errorf("RegJaegerForJsonStr needs to initialize or pass jaegerConfig parameters ")
		return
	}

	t = &Tracer{
		conf: &config,
	}
	//TODO 可以增加任意实现啦opentracing协议的链路客户端
	switch config.TracingType {
	default: //默认 jaeger
		t.tracer, t.closer, err = newJaeger(serviceName, config, op...)
	}

	if err != nil {
		return
	}

	opentracing.SetGlobalTracer(t.tracer)
	return
}

func (t *Tracer) Instance() opentracing.Tracer {
	if t.tracer == nil {
		return opentracing.GlobalTracer()
	}
	return t.tracer
}

func (t *Tracer) StartTraceOperator(r **http.Request, SpanName string, extract ...bool) (TraceOperator, error) {
	if t.tracer == nil {
		return nil, fmt.Errorf("tracer is not initialized")
	}
	return StartTraceOperator(r, SpanName, extract...)
}

func (t *Tracer) StartTraceCtx(r context.Context, SpanName string) (context.Context, TraceOperator, error) {
	if t.tracer == nil {
		return r, nil, fmt.Errorf("tracer is not initialized")
	}
	return StartTraceCtx(r, SpanName)
}

func (t *Tracer) GetTracingConfig() *TracingConfig {
	return t.conf
}

func (t *Tracer) Close() error {
	if t.closer == nil {
		return nil
	}
	return t.closer.Close()
}
