package tracer

import (
	"context"
	"testing"
	"time"
)

func TestTracing(t *testing.T) {

	//step 1：初始化tracing实例,并且通过使用AddTracer(tr)注册到http中
	_, err := New("localTest", TracingConfig{
		TracingType:  "jaeger",
		SamplerType:  "const",
		SamplerParam: "1",
		SenderType:   "udp",
		Endpoint:     "127.0.0.1:6831",
	})
	if err != nil {
		panic("parse tracing config err：" + err.Error())
	}

	//step 2：程序内存使用
	var ctx = context.Background()
	_, span, err := StartTraceCtx(ctx, "span name")
	if err != nil {
		panic("parse StartTraceOperator err：" + err.Error())
	}

	span.SetTag("key", "value")
	span.Finish()
	time.Sleep(time.Second * 2)
}
