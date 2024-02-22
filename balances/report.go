package balance

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
)

func report(ctx context.Context, action string, t time.Time) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		execDuration := time.Since(t)
		span.SetTag("balance", fmt.Sprintf("action%v  Took: %v", action, execDuration))
	}
}
