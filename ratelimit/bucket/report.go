package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
)

func report(ctx context.Context, action string, takeNum interface{}, t time.Time) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		execDuration := time.Since(t)
		span.SetTag("ratelimit", fmt.Sprintf("action%v take:%v - Took: %v", action, takeNum, execDuration))
	}
}
