package blademaster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	xhttp "net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/iooikaak/frame/conf/env"
	httpMetadata "github.com/iooikaak/frame/metadata"
	"github.com/iooikaak/frame/net/metadata"
	"github.com/iooikaak/frame/net/netutil/breaker"

	"github.com/gogo/protobuf/proto"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	pkgerr "github.com/pkg/errors"
)

const (
	_minRead     = 16 * 1024 // 16kb
	_gomsEnvName = "GOMS_OUT_HTTP_ADDR"
	_contentType = "Content-Type"
	_urlencoded  = "application/x-www-form-urlencoded"
	_json        = "application/json"
)

var (
	_noKickUserAgent = "blademaster"
	_gatewayPath     = "/v1/services"
	_gomsAddr        = ""
)

func init() {
	n, err := os.Hostname()
	if err == nil {
		_noKickUserAgent = _noKickUserAgent + runtime.Version() + " " + n
	}
	addr := os.Getenv(_gomsEnvName)
	if addr != "" {
		_gomsAddr = addr
	}
}

// ClientConfig is http client conf.
type ClientConfig struct {
	Dial      time.Duration            `yaml:"dial"`
	Timeout   time.Duration            `yaml:"timeout"`
	KeepAlive time.Duration            `yaml:"keepalive"`
	Breaker   *breaker.Config          `yaml:"breaker"`
	URL       map[string]*ClientConfig `yaml:"url"`
	Host      map[string]*ClientConfig `yaml:"host"`
}

// Client is http client.
type Client struct {
	conf      *ClientConfig
	client    *xhttp.Client
	dialer    *net.Dialer
	transport xhttp.RoundTripper

	urlConf  map[string]*ClientConfig
	hostConf map[string]*ClientConfig
	mutex    sync.RWMutex
	breaker  *breaker.Group
}

// NewClient new a http client.
func NewClient(c *ClientConfig) *Client {
	client := new(Client)
	client.conf = c
	client.dialer = &net.Dialer{
		Timeout:   c.Dial,
		KeepAlive: c.KeepAlive,
	}

	client.client = &xhttp.Client{
		Transport: &nethttp.Transport{RoundTripper: client.transport},
	}
	client.urlConf = make(map[string]*ClientConfig)
	client.hostConf = make(map[string]*ClientConfig)
	client.breaker = breaker.NewGroup(c.Breaker)
	if c.Timeout <= 0 {
		panic("must config http timeout!!!")
	}
	for uri, cfg := range c.URL {
		client.urlConf[uri] = cfg
	}
	for host, cfg := range c.Host {
		client.hostConf[host] = cfg
	}
	return client
}

// SetTransport set client transport
func (client *Client) SetTransport(t xhttp.RoundTripper) {
	client.transport = t
	client.client.Transport = &nethttp.Transport{RoundTripper: t}
}

// SetConfig set client config.
func (client *Client) SetConfig(c *ClientConfig) {
	client.mutex.Lock()
	if c.Timeout > 0 {
		client.conf.Timeout = c.Timeout
	}
	if c.KeepAlive > 0 {
		client.dialer.KeepAlive = c.KeepAlive
		client.conf.KeepAlive = c.KeepAlive
	}
	if c.Dial > 0 {
		client.dialer.Timeout = c.Dial
		client.conf.Timeout = c.Dial
	}
	if c.Breaker != nil {
		client.conf.Breaker = c.Breaker
		client.breaker.Reload(c.Breaker)
	}
	for uri, cfg := range c.URL {
		client.urlConf[uri] = cfg
	}
	for host, cfg := range c.Host {
		client.hostConf[host] = cfg
	}
	client.mutex.Unlock()
}

// NewRequest new http request with method, uri, ip, values and headers.
func (client *Client) NewRequest(method, uri, realIP string, params url.Values) (req *xhttp.Request, err error) {
	if method == xhttp.MethodGet {
		req, err = xhttp.NewRequest(xhttp.MethodGet, fmt.Sprintf("%s?%s", uri, params.Encode()), nil)
	} else {
		req, err = xhttp.NewRequest(xhttp.MethodPost, uri, strings.NewReader(params.Encode()))
	}
	if err != nil {
		err = pkgerr.Wrapf(err, "method:%s,uri:%s", method, uri)
		return
	}

	if method == xhttp.MethodPost {
		req.Header.Set(_contentType, _urlencoded)
	}
	if realIP != "" {
		req.Header.Set(_httpHeaderRemoteIP, realIP)
	}
	req.Header.Set(httpMetadata.HttpFrom, env.AppID)
	return
}

// Get issues a GET to the specified URL.
func (client *Client) Get(c context.Context, uri, ip string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodGet, uri, ip, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// Post issues a Post to the specified URL.
func (client *Client) Post(c context.Context, uri, ip string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodPost, uri, ip, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// Delete issues a Delete to the specified URL.
func (client *Client) Delete(c context.Context, uri, realIP string, params url.Values, res interface{}) (err error) {
	req, err := xhttp.NewRequest(xhttp.MethodDelete, fmt.Sprintf("%s?%s", uri, params.Encode()), nil)
	if err != nil {
		err = pkgerr.Wrapf(err, "method:delete, uri:%s", uri)
		return
	}

	if realIP != "" {
		req.Header.Set(_httpHeaderRemoteIP, realIP)
	}
	req.Header.Set(httpMetadata.HttpFrom, env.AppID)
	return client.Do(c, req, res)
}

// Put issues a Put to the specified URL.
func (client *Client) Put(c context.Context, uri, realIP string, params url.Values, res interface{}) (err error) {
	req, err := xhttp.NewRequest(xhttp.MethodPut, fmt.Sprintf("%s?%s", uri, params.Encode()), nil)
	if err != nil {
		err = pkgerr.Wrapf(err, "method:delete, uri:%s", uri)
		return
	}
	if realIP != "" {
		req.Header.Set(_httpHeaderRemoteIP, realIP)
	}
	return client.Do(c, req, res)
}

// RawPut issues a RawPut to the specified URL.
func (client *Client) RawPut(c context.Context, uri, realIP string, params url.Values, res interface{}) (resp *xhttp.Response, err error) {
	req, err := xhttp.NewRequest(xhttp.MethodPut, fmt.Sprintf("%s?%s", uri, params.Encode()), nil)
	if err != nil {
		err = pkgerr.Wrapf(err, "method:delete, uri:%s", uri)
		return
	}
	if realIP != "" {
		req.Header.Set(_httpHeaderRemoteIP, realIP)
	}
	resp, bs, err := client.RawDo(c, req)
	if err != nil {
		return
	}
	if res != nil && len(bs) > 0 {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// RESTfulGet issues a RESTful GET to the specified URL.
func (client *Client) RESTfulGet(c context.Context, uri, ip string, params url.Values, res interface{}, v ...interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodGet, fmt.Sprintf(uri, v...), ip, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res, uri)
}

// RESTfulPost issues a RESTful Post to the specified URL.
func (client *Client) RESTfulPost(c context.Context, uri, ip string, params url.Values, res interface{}, v ...interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodPost, fmt.Sprintf(uri, v...), ip, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res, uri)
}

// Raw sends an HTTP request and returns bytes response
func (client *Client) Raw(c context.Context, req *xhttp.Request, v ...string) (bs []byte, err error) {
	var (
		ok      bool
		code    string
		cancel  func()
		resp    *xhttp.Response
		config  *ClientConfig
		timeout time.Duration
		uri     = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.Path)
		tr      *nethttp.Tracer
	)
	// NOTE fix prom & config uri key.
	if len(v) == 1 {
		uri = v[0]
	}

	req.Header.Set(httpMetadata.HttpFrom, env.AppID)

	// breaker
	brk := client.breaker.Get(uri)
	if err = brk.Allow(); err != nil {
		code = "breaker"
		_metricClientReqCodeTotal.Inc(uri, req.Method, code)
		return
	}
	defer client.onBreaker(brk, &err)
	// stat
	now := time.Now()
	defer func() {
		_metricClientReqDur.Observe(int64(time.Since(now)/time.Millisecond), uri, req.Method)
		if code != "" {
			_metricClientReqCodeTotal.Inc(uri, req.Method, code)
		}
	}()
	// get config
	// 1.url config 2.host config 3.default
	client.mutex.RLock()
	if config, ok = client.urlConf[uri]; !ok {
		if config, ok = client.hostConf[req.Host]; !ok {
			config = client.conf
		}
	}
	client.mutex.RUnlock()
	// timeout
	deliver := true
	timeout = config.Timeout
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < timeout {
			// deliver small timeout
			timeout = ctimeout
			deliver = false
		}
	}
	if deliver {
		c, cancel = context.WithTimeout(c, timeout)
		defer cancel()
	}
	setTimeout(req, timeout)
	req = req.WithContext(c)
	setCaller(req)
	metadata.Range(c,
		func(key string, value interface{}) {
			setMetadata(req, key, value)
		},
		metadata.IsOutgoingKey)
	// add tracing
	if c != nil {
		// If root does not exist, it will not be automatically tracked
		if span := opentracing.SpanFromContext(c); span != nil {
			// open tracing
			req, tr = nethttp.TraceRequest(opentracing.GlobalTracer(), req, getOperationName(req.Method, req))
			defer tr.Finish()
		}
	}

	if _gomsAddr != "" && strings.HasPrefix(req.URL.Path, _gatewayPath) {
		err = directGatewayRequest(req)
		if err != nil {
			return nil, pkgerr.Wrapf(err, "method:%s,uri:%s", req.Method, req.URL.String())
		}
	}

	if resp, err = client.client.Do(req); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		code = "failed"
		return
	}
	defer func() { _ = resp.Body.Close() }()
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		return
	}
	if resp.StatusCode >= xhttp.StatusBadRequest {
		err = pkgerr.Errorf("incorrect http status:%d host:%s, url:%s body(%s)", resp.StatusCode, req.URL.Host, realURL(req), string(bs))
		code = strconv.Itoa(resp.StatusCode)
		return
	}
	return
}

// RawDo 和Raw请求类似, 但是开放了原始的response返回
func (client *Client) RawDo(c context.Context, req *xhttp.Request, v ...string) (resp *xhttp.Response, bs []byte, err error) {
	var (
		ok      bool
		code    string
		cancel  func()
		config  *ClientConfig
		timeout time.Duration
		uri     = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.Path)
		tr      *nethttp.Tracer
	)
	// NOTE fix prom & config uri key.
	if len(v) == 1 {
		uri = v[0]
	}
	req.Header.Set(httpMetadata.HttpFrom, env.AppID)

	// breaker
	brk := client.breaker.Get(uri)
	if err = brk.Allow(); err != nil {
		code = "breaker"
		_metricClientReqCodeTotal.Inc(uri, req.Method, code)
		return
	}
	defer client.onBreaker(brk, &err)
	defer func(now time.Time) {
		_metricClientReqDur.Observe(int64(time.Since(now)/time.Millisecond), uri, req.Method)
		if code != "" {
			_metricClientReqCodeTotal.Inc(uri, req.Method, code)
		}
	}(time.Now())
	// get config
	// 1.url config 2.host config 3.default
	client.mutex.RLock()
	if config, ok = client.urlConf[uri]; !ok {
		if config, ok = client.hostConf[req.Host]; !ok {
			config = client.conf
		}
	}
	client.mutex.RUnlock()
	// timeout
	deliver := true
	timeout = config.Timeout
	if deadline, ok := c.Deadline(); ok {
		if cTimeout := time.Until(deadline); cTimeout < timeout {
			// deliver small timeout
			timeout = cTimeout
			deliver = false
		}
	}
	if deliver {
		c, cancel = context.WithTimeout(c, timeout)
		defer cancel()
	}
	setTimeout(req, timeout)
	req = req.WithContext(c)
	setCaller(req)
	metadata.Range(c,
		func(key string, value interface{}) {
			setMetadata(req, key, value)
		},
		metadata.IsOutgoingKey)
	// add tracing
	if c != nil {
		// If root does not exist, it will not be automatically tracked
		if span := opentracing.SpanFromContext(c); span != nil {
			// open tracing
			req, tr = nethttp.TraceRequest(opentracing.GlobalTracer(), req, getOperationName(req.Method, req))
			defer tr.Finish()
		}
	}

	if _gomsAddr != "" && strings.HasPrefix(req.URL.Path, _gatewayPath) {
		err = directGatewayRequest(req)
		if err != nil {
			return nil, nil, pkgerr.Wrapf(err, "method:%s,uri:%s", req.Method, req.URL.String())
		}
	}
	if resp, err = client.client.Do(req); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		code = "failed"
		return
	}
	defer func() { _ = resp.Body.Close() }()
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
	}

	return
}

func getOperationName(method string, req *xhttp.Request) nethttp.ClientOption {
	return nethttp.OperationName(fmt.Sprintf("HTTP %s %s", method, req.URL.EscapedPath()))
}

// Do sends an HTTP request and returns an HTTP json response.
func (client *Client) Do(c context.Context, req *xhttp.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil && len(bs) > 0 {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// JSON sends an HTTP request and returns an HTTP json response.
func (client *Client) JSON(c context.Context, req *xhttp.Request, res interface{}, v ...string) (err error) {
	req.Header.Set(_contentType, _json)
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil && len(bs) > 0 {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// PB sends an HTTP request and returns an HTTP proto response.
func (client *Client) PB(c context.Context, req *xhttp.Request, res proto.Message, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil && len(bs) > 0 {
		if err = proto.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// Post issues a Post to the specified URL.
func (client *Client) PostJson(c context.Context, uri, ip string, data interface{}, res interface{}) (err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := xhttp.NewRequest(xhttp.MethodPost, uri, bytes.NewReader(b))
	if err != nil {
		return
	}
	if ip != "" {
		req.Header.Set(_httpHeaderRemoteIP, ip)
	}
	req.Header.Set(_contentType, _json)
	return client.Do(c, req, res)
}

func (client *Client) onBreaker(breaker breaker.Breaker, err *error) {
	if err != nil && *err != nil {
		breaker.MarkFailed()
	} else {
		breaker.MarkSuccess()
	}
}

// realUrl return url with http://host/params.
func realURL(req *xhttp.Request) string {
	if req.Method == xhttp.MethodGet {
		return req.URL.String()
	} else if req.Method == xhttp.MethodPost {
		ru := req.URL.Path
		if req.Body != nil {
			rd, ok := req.Body.(io.Reader)
			if ok {
				buf := bytes.NewBuffer([]byte{})
				_, _ = buf.ReadFrom(rd)
				ru = ru + "?" + buf.String()
			}
		}
		return ru
	}
	return req.URL.Path
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}

func directGatewayRequest(req *xhttp.Request) error {
	var (
		actions     []string
		serviceName string // 转化后服务名
		action      string
		q           = req.URL.Query()
	)
	if q.Get("crossTheCluster") == "force" {
		return nil
	}

	actions = strings.Split(strings.Replace(req.URL.Path, _gatewayPath, "", -1), "/")
	if len(actions) < 3 {
		return pkgerr.New("请求服务网关, 但是方法不符合规范")
	}
	serviceName = strings.Replace(strings.Replace(actions[1], "-", ".", 1), "-", "", -1)
	action = fmt.Sprintf("%s.%s", serviceName, strings.Join(actions[2:], "."))

	// parse query
	req.URL.Host = _gomsAddr
	req.URL.Path = "/"
	req.URL.Scheme = "http"

	q.Set("goms_action", action)
	req.URL.RawQuery = q.Encode()
	return nil
}
