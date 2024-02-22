package elasticsearch

import (
	"bytes"
	"fmt"

	//"github.com/iooikaak/frame/elastic"
	"io/ioutil"
	"net"
	xhttp "net/http"
	"strings"
	"time"

	"github.com/iooikaak/frame/elastic"

	v7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
)

type ElasticConfig struct {
	Addrs              []string      `yaml:"addrs" json:"addrs"`
	Username           string        `yaml:"username" json:"username"`
	Password           string        `yaml:"password" json:"password"`
	HealthcheckEnabled bool          `yaml:"healthcheckEnabled" json:"healthcheckEnabled"`
	SnifferEnabled     bool          `yaml:"snifferEnabled" json:"snifferEnabled"`
	HealthTimeOut      time.Duration `yaml:"healthtimeout" json:"healthtimeout"`
	SnifferTimeout     time.Duration `yaml:"snifferTimeout" json:"snifferTimeout"`
	V7                 struct {
		MaxIdleConnsPerHost int `yaml:"maxIdleConnsPerHost" json:"maxIdleConnsPerHost"`
		MaxIdleConns        int `yaml:"maxIdleConns" json:"maxIdleConns"`
		TimeOut             int `yaml:"timeOut" json:"timeOut"`
		KeepAlive           int `yaml:"keepAlive" json:"keepAlive"`
	} `yaml:"v7" json:"v7"`
}

const (
	es7                = 7
	clientTimeout      = 30
	clientKeepAlive    = 30
	clientMaxIdleConns = 100
)

func New(esConfig *ElasticConfig, option ...elastic.ClientOptionFunc) (es *elastic.Client, err error) {

	if len(esConfig.Addrs) == 0 {
		err = fmt.Errorf("addrs is empty")
		return
	}

	if esConfig.HealthTimeOut == 0 {
		esConfig.HealthTimeOut = elastic.DefaultHealthcheckTimeout
	}

	if esConfig.SnifferTimeout == 0 {
		esConfig.SnifferTimeout = elastic.DefaultSnifferTimeout
	}

	if esConfig.SnifferTimeout == 0 {
		esConfig.SnifferTimeout = elastic.DefaultSnifferTimeout
	}
	//提供些默认必要配置，但还支持elastic更多配置，提供啦注入配置方式
	option = append(option, elastic.SetURL(esConfig.Addrs...),
		elastic.SetBasicAuth(esConfig.Username, esConfig.Password),
		elastic.SetHealthcheck(esConfig.HealthcheckEnabled),
		elastic.SetHealthcheckTimeout(esConfig.HealthTimeOut),
		elastic.SetSnifferTimeout(esConfig.SnifferTimeout),
		elastic.SetSniff(esConfig.SnifferEnabled))
	es, err = elastic.NewClient(option...)

	return
}

func NewV7(esConfig *ElasticConfig, option ...v7.ClientOptionFunc) (es *v7.Client, err error) {

	if len(esConfig.Addrs) == 0 {
		err = fmt.Errorf("addrs is empty")
		return
	}

	var (
		timeout             = esConfig.V7.TimeOut
		keepAlive           = esConfig.V7.KeepAlive
		maxIdleConns        = esConfig.V7.MaxIdleConns
		maxIdleConnsPerHost = esConfig.V7.MaxIdleConnsPerHost
	)

	if esConfig.HealthTimeOut == 0 {
		esConfig.HealthTimeOut = v7.DefaultHealthcheckTimeout
	}

	if esConfig.SnifferTimeout == 0 {
		esConfig.SnifferTimeout = v7.DefaultSnifferTimeout
	}

	if esConfig.SnifferTimeout == 0 {
		esConfig.SnifferTimeout = v7.DefaultSnifferTimeout
	}

	if timeout == 0 {
		timeout = clientTimeout
	}

	if keepAlive == 0 {
		keepAlive = clientKeepAlive
	}

	if maxIdleConns == 0 {
		maxIdleConns = clientMaxIdleConns
	}

	if maxIdleConnsPerHost == 0 {
		maxIdleConnsPerHost = xhttp.DefaultMaxIdleConnsPerHost
	}
	if esConfig.V7.MaxIdleConns == 0 {
		esConfig.SnifferTimeout = v7.DefaultSnifferTimeout
	}

	// 默认官方http配置
	transport := &xhttp.Transport{
		Proxy: xhttp.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			KeepAlive: time.Duration(keepAlive) * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          maxIdleConns,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   maxIdleConnsPerHost,
	}

	// 提供些默认必要配置，但还支持elastic更多配置，提供啦注入配置方式
	option = append(option, v7.SetURL(esConfig.Addrs...),
		v7.SetHttpClient(&TracingClient{xhttp.Client{Transport: transport}}),
		v7.SetBasicAuth(esConfig.Username, esConfig.Password),
		v7.SetHealthcheck(esConfig.HealthcheckEnabled),
		v7.SetHealthcheckTimeout(esConfig.HealthTimeOut),
		v7.SetSnifferTimeout(esConfig.SnifferTimeout),
		v7.SetSniff(esConfig.SnifferEnabled))
	es, err = v7.NewClient(option...)
	return
}

type TracingClient struct {
	xhttp.Client
}

func (t *TracingClient) Do(req *xhttp.Request) (*xhttp.Response, error) {
	var (
		resp      *xhttp.Response
		err       error
		bodyBytes []byte
	)

	if req.Body != nil && strings.Contains(req.URL.Path, "_search") {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	defer func(expireTime time.Time, span opentracing.Span) {
		t.report(span, req, resp, err, expireTime, bodyBytes)
	}(time.Now(), t.startSpan(req))
	resp, err = t.Client.Do(req)
	return resp, err
}

func (t *TracingClient) startSpan(req *xhttp.Request) (span opentracing.Span) {
	if span = opentracing.SpanFromContext(req.Context()); span != nil {
		span = opentracing.StartSpan("elastic-"+req.Method+":"+req.URL.Path, opentracing.ChildOf(span.Context()))
	}

	return
}

func (t *TracingClient) report(span opentracing.Span, req *xhttp.Request, res *xhttp.Response, err error, tt time.Time, bodyBytes []byte) {

	var (
		responseStateCode int
	)

	if span != nil {
		defer span.Finish()

		if res != nil {
			responseStateCode = res.StatusCode
		}

		execDuration := time.Since(tt)
		span.SetTag(req.Method, req.URL.Path)
		span.SetTag(req.Method+":"+req.URL.Path, fmt.Sprintf("condition:%+v responseCode:%d err:%v   - Took: %v", string(bodyBytes), responseStateCode, err, execDuration))
	}
}
