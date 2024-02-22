
# 快速开始
* metrics使用例子

```go

type WriteMetricsVet struct {
	Attempts   metrics.Counter `metric:"attempts" tags:"action=retry" help:"Number of retries"`
	Inserts    metrics.Counter `metric:"inserts" tags:"action=success" help:"The number of times from mass production to es success"`
	Errors     metrics.Counter `metric:"errors" tags:"action=error" help:"Number of times from mass production to es error"`
	LatencyOk  metrics.Timer   `metric:"latency-ok" tags:"action=LatencyOk" help:"Time from mass production to es into power consumption"`
	LatencyErr metrics.Timer   `metric:"latency-err" tags:"action=LatencyErr" help:"Time from mass production to es failure"`
}
          //1. 实例化一个metrics
            // new metrics
            met := metrics.New('serverName')
            //2. 注入需要的metric数据,WriteMetrics是业务中需要统计的指标聚合信息
                t := &WriteMetrics{}
            	m.Register(t, nil)

          //业务中使用
          请看go-example monitor例子
```
