package tracer

import (
	"context"
	"io"
	"net/http"

	"github.com/iooikaak/frame/config/env"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

func newJaeger(serviceName string, tracingConfig TracingConfig, op ...Option) (opentracing.Tracer, io.Closer, error) {
	var (
		sampler  jaeger.Sampler
		reporter jaeger.Reporter
		config   = tracingConfig
		err      error
	)
	opts := applyOptions(op...)
	sampler, reporter, err = newSamplerAndReporter(config, opts)
	if err != nil {
		return nil, nil, err
	}
	tracing, cl := jaeger.NewTracer(
		serviceName,
		sampler,
		reporter,
		jaeger.TracerOptions.Metrics(opts.metricInstance),
		jaeger.TracerOptions.Logger(opts.logger),
		jaeger.TracerOptions.MaxTagValueLength(opts.maxTagValueLength),
		jaeger.TracerOptions.PoolSpans(opts.poolSpans),
		jaeger.TracerOptions.Tag("env", env.DeployEnv),
	)

	tracingType = jaegerTracing
	return tracing, cl, nil
}

type jaegerTracingOperator struct {
	span      opentracing.Span
	tracingId string
	spanId    string
	parentId  string
}

func startJaegerTracingOperator(r **http.Request, SpanName string, extract bool) (t TraceOperator, err error) {

	var (
		span   opentracing.Span
		option opentracing.StartSpanOption = opentracing.FollowsFrom(nil)
	)

	if r != nil {
		if extract {
			carrier := opentracing.HTTPHeadersCarrier((*r).Header)
			ctx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
			option = ext.RPCServerOption(ctx)
		} else {
			s := opentracing.SpanFromContext((*r).Context())
			if s != nil {
				option = opentracing.ChildOf(s.Context())
			}
		}
	}

	span = opentracing.StartSpan(SpanName, option)

	ctx := span.Context().(jaeger.SpanContext)

	t = &jaegerTracingOperator{
		span:      span,
		tracingId: ctx.TraceID().String(),
		spanId:    ctx.SpanID().String(),
		parentId:  ctx.ParentID().String(),
	}

	if r != nil {
		*r = (*r).WithContext(opentracing.ContextWithSpan((*r).Context(), span))
	}

	return
}

func startTracingCtx(r context.Context, SpanName string) (c context.Context, t TraceOperator, err error) {

	var (
		span   opentracing.Span
		option opentracing.StartSpanOption = opentracing.FollowsFrom(nil)
	)
	if r != nil {
		s := opentracing.SpanFromContext(r)
		if s != nil {
			option = opentracing.ChildOf(s.Context())
		}
	}

	span = opentracing.StartSpan(SpanName, option)
	ctx := span.Context().(jaeger.SpanContext)
	t = &jaegerTracingOperator{
		span:      span,
		tracingId: ctx.TraceID().String(),
		spanId:    ctx.SpanID().String(),
		parentId:  ctx.ParentID().String(),
	}

	if r != nil {
		c = opentracing.ContextWithSpan(r, span)
	} else {
		c = opentracing.ContextWithSpan(context.Background(), span)
	}

	return
}

//数据会跟着链路走
func (t *jaegerTracingOperator) SetTracingData(k, v string) {
	t.span.SetBaggageItem(k, v)
}

func (t *jaegerTracingOperator) GetTracingData(k string) string {
	return t.span.BaggageItem(k)
}

//局部span数据，可做为后台搜索关键字
func (t *jaegerTracingOperator) SetTag(k string, v interface{}) {
	t.span.SetTag(k, v)
}

func (t *jaegerTracingOperator) Inject(r *http.Request) error {
	return t.span.Tracer().Inject(t.span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
}

func (t *jaegerTracingOperator) Extract(header http.Header) (opentracing.SpanContext, error) {
	return t.span.Tracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
}

func (t *jaegerTracingOperator) GetTracingId() string {
	return t.tracingId
}

func (t *jaegerTracingOperator) GetSpanId() string {
	return t.spanId
}

func (t *jaegerTracingOperator) GetParentId() string {
	return t.parentId
}

//提供span
func (t *jaegerTracingOperator) GetSpan() opentracing.Span {
	return t.span
}

func (t *jaegerTracingOperator) Finish() {
	t.span.Finish()
}
