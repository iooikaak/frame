package gins

import (
	"fmt"
	"net/url"

	"github.com/iooikaak/frame/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

const _reqID = "shreqid"

func NewGinJaegerTrace(tr opentracing.Tracer) HandlerFunc {

	return func(c *Context) {
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		ctx, _ := tr.Extract(opentracing.HTTPHeaders, carrier)
		urlStr := getEscapedPath(c.Request.URL)
		op := fmt.Sprintf("HTTP %s %s", c.Request.Method, urlStr)
		sp := tr.StartSpan(op, ext.RPCServerOption(ctx))
		ext.HTTPMethod.Set(sp, c.Request.Method)
		ext.HTTPUrl.Set(sp, urlStr)
		componentName := "net/http"
		ext.Component.Set(sp, componentName)

		c.C = opentracing.ContextWithSpan(c.C, sp)
		// opentracing协议没有暴露getTraceId方法, 所以获取比较暴力
		if jaegerCtx, ok := sp.Context().(jaeger.SpanContext); ok {
			c.Set(metadata.HttpTraceId, jaegerCtx.TraceID().String())
		}
		c.Next()

		ext.HTTPStatusCode.Set(sp, uint16(c.Writer.Status()))
		sp.Finish()
	}
}

func getEscapedPath(u *url.URL) string {
	url := u.EscapedPath()
	if url == "" {
		url = "/"
	}

	return url
}
