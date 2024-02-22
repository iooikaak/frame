package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uber/jaeger-lib/metrics"
	pmetrics "github.com/uber/jaeger-lib/metrics/prometheus"
)

var (
	defaultMetricPrefix = "micro"
	defaultMetric       *Metrics
)

type Metrics struct {
	reg            prometheus.Registerer
	gather         prometheus.Gatherer
	metricsFactory metrics.Factory
	globalTags     map[string]string
}

func New(serverName string) *Metrics {

	globalTags := make(map[string]string)
	if len(serverName) > 0 {
		globalTags[fmt.Sprintf("%s_%s", defaultMetricPrefix, "name")] = serverName
	}

	prometheus.DefaultRegisterer = prometheus.WrapRegistererWith(globalTags, prometheus.DefaultRegisterer)
	defaultMetric = &Metrics{
		reg:            prometheus.DefaultRegisterer,
		gather:         prometheus.DefaultGatherer,
		globalTags:     globalTags,
		metricsFactory: pmetrics.New(pmetrics.WithRegisterer(prometheus.DefaultRegisterer)),
	}

	return defaultMetric
}

func (m Metrics) Instance() metrics.Factory {
	return m.metricsFactory
}

func (m Metrics) RegInstance() prometheus.Registerer {
	return m.reg
}

func (m Metrics) Gather() prometheus.Gatherer {
	return m.gather
}

func (m *Metrics) Register(i interface{}, tags map[string]string) {
	if tags != nil {
		for k, v := range m.globalTags {
			tags[k] = v
		}
	}
	metrics.MustInit(i, m.metricsFactory, tags)
}

func Register(i interface{}, tags map[string]string) error {
	if defaultMetric == nil {
		return fmt.Errorf("must init metrics.New")
	}

	defaultMetric.Register(i, tags)
	return nil
}

func Instance() metrics.Factory {
	if defaultMetric == nil {
		return nil
	}
	return defaultMetric.Instance()
}
