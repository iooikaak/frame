package gin

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/url"
)

func (engine *Engine) AddTracer(tr opentracing.Tracer) {
	if tr != nil {
		engine.Use(NewGinJaegerTrace(tr))
	}
}

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
		c.Request = c.Request.WithContext(
			opentracing.ContextWithSpan(c.Request.Context(), sp))

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
