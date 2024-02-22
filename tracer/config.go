package tracer

import (
	"fmt"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"net/url"
	"strconv"
	"strings"
)

type TracingConfig struct {
	TracingType  string `yaml:"tracingType" json:"tracingType"`
	SamplerType  string `yaml:"samplerType" json:"samplerType"`
	SamplerParam string `yaml:"samplerParam" json:"samplerParam"`
	SenderType   string `yaml:"senderType" json:"senderType"`
	Endpoint     string `yaml:"endpoint" json:"endpoint"`
}

func newSamplerAndReporter(config TracingConfig, ops Options) (jaeger.Sampler, jaeger.Reporter, error) {

	var (
		sampler      jaeger.Sampler
		reporter     jaeger.Reporter
		tempSampler  jaeger.Sampler
		tempReporter jaeger.Reporter
		err          error
	)

	tempSampler, err = newSampler(config.SamplerType, config.SamplerParam)
	if err == nil {
		sampler = tempSampler
	}

	tempReporter, err = newReporter(config.SenderType, config.Endpoint, ops)
	if err == nil {
		reporter = tempReporter
	}

	if sampler == nil {
		return nil, nil, fmt.Errorf("sampler Unavailable Jaeger sampler configuration err [%v] ", err)
	}

	if reporter == nil {
		return nil, nil, fmt.Errorf("reporter Unavailable Jaeger reporter configuration err [%v]", err)
	}

	return sampler, reporter, nil
}

//Sampler
func newSampler(typ string, param string) (sampler jaeger.Sampler, err error) {
	if typ == "" || param == "" {
		err = fmt.Errorf("Unavailable Jaeger sampler configuration")
		return
	}

	switch strings.ToLower(typ) {
	case "const":
		sampler = jaeger.NewConstSampler(param == "1")
	case "probabilistic":
		var samplingRate float64
		samplingRate, err = strconv.ParseFloat(param, 64)
		if err != nil {
			return nil, err
		}
		sampler, err = jaeger.NewProbabilisticSampler(samplingRate)
		if err != nil {
			return nil, err
		}
	case "ratelimiting":
		var maxTracesPerSecond float64
		maxTracesPerSecond, err = strconv.ParseFloat(param, 64)
		if err != nil {
			return nil, err
		}
		sampler = jaeger.NewRateLimitingSampler(maxTracesPerSecond)
	case "remote":
		sampler = jaeger.NewRemotelyControlledSampler(param)
	default:
		err = fmt.Errorf("Jaeger does not have this type(%s) of sampler", typ)
		return
	}

	return sampler, nil
}

//Reporter
func newReporter(senderType string, endpoint string, ops Options) (reporter jaeger.Reporter, err error) {
	if senderType == "" || endpoint == "" {
		err = fmt.Errorf("Unavailable Jaeger reporter configuration")
		return
	}

	var sender jaeger.Transport

	switch strings.ToLower(senderType) {
	case "udp":
		sender, err = jaeger.NewUDPTransport(endpoint, 60000)
		if err != nil {
			return
		}
	case "http":
		sender = transport.NewHTTPTransport(endpoint)
	case "kafka":
		sender, err = NewKafkaTransport(endpoint)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("Jaeger does not have this type(%s) of sender", senderType)
		return
	}

	reporter = jaeger.NewRemoteReporter(sender,
		jaeger.ReporterOptions.Metrics(ops.metricInstance),
		jaeger.ReporterOptions.Logger(ops.logger),
		jaeger.ReporterOptions.QueueSize(ops.queueSize))
	return
}

func GetEscapedPath(u *url.URL) string {
	url := u.EscapedPath()
	if url == "" {
		url = "/"
	}

	return url
}
