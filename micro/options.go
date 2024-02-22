package micro

import (
	"github.com/micro/go-micro/v2"
	"net"
	"time"

	"github.com/micro/go-micro/v2/server"

	"github.com/iooikaak/frame/stat/metric/metrics"
	"github.com/iooikaak/frame/xlog"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/transport"
	"github.com/opentracing/opentracing-go"
)

const (
	microRegister    = "consul"
	microRpc         = "grpc"
	defaultMicroPort = ":10086"
)

type Option func(c *Options)

type Options struct {
	// base info
	ServerName    string `yaml:"serverName" json:"serverName"`
	User          string `yaml:"user" json:"user"`
	Password      string `yaml:"password" json:"password"`
	ServerVersion string `yaml:"serverVersion" json:"serverVersion"`
	Address       string `yaml:"address" json:"address"`

	// client
	PoolSize       int           `yaml:"poolSize" json:"poolSize"`
	Retries        int           `yaml:"retries" json:"retries"`
	PoolTTL        time.Duration `yaml:"poolTTL" json:"poolTTL"`
	RequestTimeout time.Duration `yaml:"requestTimeout" json:"requestTimeout"`
	DialTimeout    time.Duration `yaml:"dialTimeout" json:"dialTimeout"`

	// rpc register
	RpcRegister string `yaml:"rpcRegister" json:"rpcRegister"`
	// rpc type
	RpcModel string `yaml:"rpcModel" json:"rpcModel"`
	// register
	CenterAddr                     []string `yaml:"centerAddr" json:"centerAddr"`                                         // 注册中心地址
	CenterName                     string   `yaml:"centerName" json:"centerName"`                                         // 注册中心注册自己的名称
	RegisterTTL                    int      `yaml:"registerTTL" json:"registerTTL"`                                       // 注册服务的有效时间 单位s 默认30s
	RegisterInterval               int      `yaml:"registerInterval" json:"registerInterval"`                             // 重新注册时间间隔（心跳）单位s 默认10s
	RenewalIntervalInSecs          int      `yaml:"renewalIntervalInSecs" json:"renewalIntervalInSecs"`                   // 重新获取服务列表时间间隔， 单位s，默认30s
	RollDiscoveriesIntervalSeconds int      `yaml:"rollDiscoveriesIntervalSeconds" json:"rollDiscoveriesIntervalSeconds"` // 滚动服务发现地址(集群环境)，单位s, 默认60s
	StatusUrl                      string   `yaml:"statusUrl" json:"statusUrl"`                                           // status url
	HeathUrl                       string   `yaml:"heathUrl" json:"heathUrl"`                                             // 健康检查url
	ServerOnly                     bool     `yaml:"serverOnly" json:"serverOnly"`                                         // 是否只作为服务端

	// log
	logLevel LogLevel
	log      xlog.ILog

	// tracing
	tracer opentracing.Tracer
	// close tracer
	DisableTracer bool `yaml:"DisableTracer" json:"DisableTracer"`
	// enable breaker
	EnableBreaker bool `yaml:"enableBreaker" json:"enableBreaker"`

	// metrics
	metrics *metrics.Metrics

	// micro.service
	l           net.Listener
	serviceFunc func(micro.Service)

	metadata map[string]string

	// micro.service before and after
	beforeStart []func()
	beforeStop  []func()
	afterStart  []func()
	afterStop   []func()

	// Wrapper
	serverWrapper []server.Option
	clientWrapper []client.Option
}

func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.metadata = md
	}
}

func Metrics(m *metrics.Metrics) Option {
	return func(c *Options) {
		c.metrics = m
	}
}

func RpcModel(rpc string) Option {
	return func(c *Options) {
		c.RpcModel = rpc
	}
}

func RegCenterAddr(addr []string) Option {
	return func(c *Options) {
		c.CenterAddr = addr
	}
}

func PoolSize(s int) Option {
	return func(c *Options) {
		c.PoolSize = s
	}
}

func PoolTTL(ttl time.Duration) Option {
	return func(c *Options) {
		c.PoolTTL = ttl
	}
}

func Retries(r int) Option {
	return func(c *Options) {
		c.Retries = r
	}
}

func RequestTimeout(timeout time.Duration) Option {
	return func(c *Options) {
		c.RequestTimeout = timeout
	}
}

func DialTimeout(ttl time.Duration) Option {
	return func(c *Options) {
		c.DialTimeout = ttl
	}
}

func Tracer(t opentracing.Tracer) Option {
	return func(c *Options) {
		c.tracer = t
	}
}

func Logger(l xlog.ILog) Option {
	return func(c *Options) {
		c.log = l
	}
}

func ServerVersion(ver string) Option {
	return func(c *Options) {
		c.ServerVersion = ver
	}
}

func ServerName(name string) Option {
	return func(c *Options) {
		c.ServerName = name
	}
}

func RegisterInterval(t int) Option {
	return func(c *Options) {
		c.RegisterInterval = t
	}
}

func RegisterTTL(t int) Option {
	return func(c *Options) {
		c.RegisterTTL = t
	}
}

func RenewalIntervalInSecs(t int) Option {
	return func(c *Options) {
		c.RenewalIntervalInSecs = t
	}
}

func RollDiscoveriesIntervalSeconds(t int) Option {
	return func(c *Options) {
		c.RollDiscoveriesIntervalSeconds = t
	}
}

func DataCenterName(dataCenterName string) Option {
	return func(c *Options) {
		c.CenterName = dataCenterName
	}
}

func Address(addr string) Option {
	return func(c *Options) {
		c.Address = addr
	}
}

func LoggerLevel(level LogLevel) Option {
	return func(c *Options) {
		c.logLevel = level
	}
}

func MicroService(f func(s micro.Service)) Option {
	return func(a *Options) {
		a.serviceFunc = f
	}
}

func BeforeStart(fn func()) Option {
	return func(a *Options) {
		a.beforeStart = append(a.beforeStart, fn)
	}
}

func BeforeStop(fn func()) Option {
	return func(a *Options) {
		a.beforeStop = append(a.beforeStop, fn)
	}
}

func AfterStart(fn func()) Option {
	return func(a *Options) {
		a.afterStart = append(a.afterStart, fn)
	}
}

func AfterStop(fn func()) Option {
	return func(a *Options) {
		a.afterStop = append(a.afterStop, fn)
	}
}

func Listener(l net.Listener) Option {
	return func(a *Options) {
		a.l = l
	}
}

// DisableTrace
func DisableTrace(t bool) Option {
	return func(a *Options) {
		a.DisableTracer = t
	}
}

// Breaker 启动断路器
func Breaker(t bool) Option {
	return func(a *Options) {
		a.EnableBreaker = t
	}
}

// ServerWrapper 服务端拦截器
func ServerWrapper(s []server.Option) Option {
	return func(a *Options) {
		a.serverWrapper = s
	}
}

// ClientWrapper 客户端拦截器
func ClientWrapper(c []client.Option) Option {
	return func(a *Options) {
		a.clientWrapper = c
	}
}

func (o *Options) Breaker() bool {
	return o.EnableBreaker
}

func (o *Options) DisableTrace() bool {
	return o.DisableTracer
}

func RpcRegister(register string) Option {
	return func(c *Options) {
		c.RpcRegister = register
	}
}

func HeathUrl(addr string) Option {
	return func(c *Options) {
		c.HeathUrl = addr
	}
}

func StatusUrl(addr string) Option {
	return func(c *Options) {
		c.StatusUrl = addr
	}
}

func applyOptions(opts *Options, options ...Option) *Options {
	for _, option := range options {
		option(opts)
	}

	if opts.logLevel == 0 {
		opts.logLevel = InfoLevel
	}

	if opts.log == nil {
		opts.log = xlog.Logger()
	}

	if len(opts.RpcModel) == 0 {
		opts.RpcModel = microRpc
	}

	if opts.PoolSize <= 0 {
		opts.PoolSize = client.DefaultPoolSize
	}

	if opts.PoolTTL <= 0 {
		opts.PoolTTL = client.DefaultPoolTTL
	}

	if opts.Retries <= 0 {
		opts.Retries = client.DefaultRetries
	}

	if opts.RegisterTTL <= 0 {
		opts.RegisterTTL = 30
	}

	if opts.tracer == nil {
		opts.tracer = opentracing.GlobalTracer()
	}

	if opts.RegisterInterval <= 0 {
		opts.RegisterInterval = 10
	}

	if len(opts.ServerVersion) == 0 {
		opts.ServerVersion = "1.0.0"
	}

	if opts.RequestTimeout <= 0 {
		opts.RequestTimeout = client.DefaultRequestTimeout
	}

	if opts.DialTimeout <= 0 {
		opts.DialTimeout = transport.DefaultDialTimeout
	}

	if len(opts.Address) == 0 {
		opts.Address = defaultMicroPort
	}

	if len(opts.RpcRegister) == 0 {
		opts.RpcRegister = microRegister
	}

	if opts.RenewalIntervalInSecs <= 0 {
		opts.RenewalIntervalInSecs = 30
	}

	if opts.RollDiscoveriesIntervalSeconds <= 0 {
		opts.RollDiscoveriesIntervalSeconds = 60
	}

	if len(opts.CenterName) == 0 {
		opts.CenterName = "discovery"
	}

	return opts
}
