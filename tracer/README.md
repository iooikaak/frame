# 项目简介

1.以jaeger为实现层的opentracing链路追踪

# 快速开始

* tracing使用例子

```go
//step 1：初始化tracing实例,并且通过使用AddTracer(tr)注册到http中
_, err := New("tracingTest", TracingConfig{
    TracingType:  "jaeger",
    SamplerType:  "const",
    SamplerParam: "1",
    SenderType:   "udp",
    Endpoint:     "10.180.18.20:6831",
})
if err != nil {
    panic("parse tracing config err：" + err.Error())
}

//step 2：程序内存使用
var ctx = context.Background()
ctx, span, err := StartTraceCtx(ctx, "span name")
if err != nil {
    panic("parse StartTraceOperator err：" + err.Error())
}

span.SetTag("key", "value")
span.Finish()

```
