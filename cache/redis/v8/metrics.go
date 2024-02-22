package redis

import (
	"strings"
	"time"

	"github.com/iooikaak/frame/stat/metric"

	"github.com/go-redis/redis/v8"
)

const namespace = "redis_v8_client"

var (
	_metricReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "redis client requests duration(ms).",
		Labels:    []string{"name", "addr", "command"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricReqErr = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "error_total",
		Help:      "redis client requests error count.",
		Labels:    []string{"name", "addr", "command", "error"},
	})
	_metricConnTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "total",
		Help:      "redis client connections total count.",
		Labels:    []string{"name", "addr", "state"},
	})
	_metricConnCurrent = metric.NewGaugeVec(&metric.GaugeVecOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "current",
		Help:      "redis client connections current.",
		Labels:    []string{"name", "addr", "state"},
	})
	_metricHits = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "",
		Name:      "hits_total",
		Help:      "redis client hits total.",
		Labels:    []string{"name", "addr"},
	})
	_metricMisses = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "",
		Name:      "misses_total",
		Help:      "redis client misses total.",
		Labels:    []string{"name", "addr"},
	})
)

func (o OpenTracingHook) report(pipe bool, elapsed time.Duration, cmds ...redis.Cmder) {
	address := o.cfg.Addr
	name := o.cfg.Name
	errStr := ""
	cmdStr := ""
	//pipeStr := fmt.Sprintf("%t", pipe)

	for _, cmd := range cmds {
		cmdStr += cmd.Name() + ";"

		if err := cmd.Err(); err != nil && err != redis.Nil {
			errStr += err.Error() + ";"
		}

		if cmd.Err() == redis.Nil {
			_metricMisses.Inc(name, address) //未命中
		}

		if cmd.Err() == nil {
			_metricHits.Inc(name, address) //命中缓存
		}
	}
	cmdStr = strings.TrimSuffix(cmdStr, ";")

	if len(errStr) > 0 {
		_metricReqErr.Inc(name, address, cmdStr, errStr)
	}

	_metricReqDur.Observe(int64(elapsed.Seconds()), name, address, cmdStr)

	_metricConnTotal.Add(float64(o.status.Hits), name, address, "total")

	_metricConnCurrent.Set(float64(o.status.TotalConns), name, address, "total")
}
