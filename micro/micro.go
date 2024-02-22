package micro

import (
	eurekaOpt "github.com/iooikaak/frame/eureka"
	"github.com/iooikaak/frame/micro/internal/breaker/hystrix"
	"github.com/iooikaak/frame/micro/internal/consul"
	"github.com/iooikaak/frame/micro/internal/etcd"
	"github.com/iooikaak/frame/micro/internal/eureka"
	"github.com/iooikaak/frame/micro/internal/metrics"
	"github.com/iooikaak/frame/micro/internal/tracer"
	"github.com/iooikaak/frame/xlog"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/config/cmd"
	mlog "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	_consulAddr   = []string{"127.0.0.1:8500"}
	_etcdAddr     = []string{"127.0.0.1:2379"}
	_etcdUser     = "root"
	_etcdPassword = ""
	_eurekaAddr   = []string{"127.0.0.1:8761"}
)

func init() {
	// 通过环境变量控制consul地址
	r := strings.TrimSpace(os.Getenv("CONSUL_IP"))
	if r != "" {
		_consulAddr = strings.Split(r, ",")
	}
	// 通过环境变量控制eureka地址
	r = strings.TrimSpace(os.Getenv("EUREKA_IP"))
	if r != "" {
		_eurekaAddr = strings.Split(r, ",")
	}
}

type Service struct {
	micro.Service
	opt      *Options
	register registry.Registry
}

func New(options *Options, o ...Option) *Service {
	name := options.ServerName
	if len(name) == 0 {
		panic("ServerName is empty")
	}

	// If micro is the suffix, it will not be appended
	if !strings.HasSuffix(name, "micro") {
		name = name + ".micro"
	}

	service := &Service{
		opt: applyOptions(options, append(o, ServerName(name))...),
	}

	return service
}

func (s *Service) Option() *Options {
	return s.opt
}

func (s *Service) Initialize(o ...Option) (err error) {
	// overflow option
	if len(o) > 0 {
		for _, v := range o {
			v(s.opt)
		}
	}

	switch s.opt.RpcRegister {
	case "eureka":
		centerAddr := _eurekaAddr
		if len(s.opt.CenterAddr) > 0 {
			centerAddr = s.opt.CenterAddr
		}
		if len(centerAddr) == 0 {
			panic("env EUREKA_IP is empty")
		}
		portS := strings.Split(s.opt.Address, ":")
		port, _ := strconv.Atoi(portS[1])
		s.register = eureka.NewRegistry(&eurekaOpt.Config{
			DefaultZone:                    centerAddr,
			RenewalIntervalInSecs:          s.opt.RenewalIntervalInSecs,
			RegistryFetchIntervalSeconds:   s.opt.RegisterInterval,
			RollDiscoveriesIntervalSeconds: s.opt.RollDiscoveriesIntervalSeconds,
			DurationInSecs:                 s.opt.RegisterTTL,
			DataCenterName:                 s.opt.CenterName,
			App:                            s.opt.ServerName,
			Port:                           port,
			StatusUrl:                      s.opt.StatusUrl,
			HealthUrl:                      s.opt.HeathUrl,
			ServerOnly:                     s.opt.ServerOnly,
		})
	case "consul":
	case "etcd":
		etcdAddr := _etcdAddr
		user := _etcdUser
		password := _etcdPassword
		if len(s.opt.CenterAddr) > 0 {
			etcdAddr = s.opt.CenterAddr
		}
		if s.opt.User != "" {
			user = s.opt.User
		}
		if s.opt.Password != "" {
			password = s.opt.Password
		}
		if len(etcdAddr) == 0 {
			panic("env ETCD_IP is empty")
		}
		s.register = etcd.NewRegistry(registry.Addrs(etcdAddr...), etcd.Auth(user, password))
	default:
		centerAddr := _consulAddr
		if len(s.opt.CenterAddr) > 0 {
			centerAddr = s.opt.CenterAddr
		}
		if len(centerAddr) == 0 {
			panic("env CONSUL_IP is empty")
		}
		s.register = consul.NewRegistry(registry.Addrs(centerAddr...))
	}

	// param init
	err = s.initMisc()
	if err != nil {
		return
	}

	// server init
	err = s.initServer()
	if err != nil {
		return
	}

	// client init
	err = s.initClient()

	return
}

func (s *Service) initMisc() (err error) {
	var (
		opts []micro.Option
	)

	// micro log init
	mlog.DefaultLogger = mlog.NewHelper(NewMicroLogger(s.opt.log, mlog.WithLevel(microLogLevel[s.opt.logLevel])))
	// auth rpc type
	if _, ok := cmd.DefaultServers[s.opt.RpcModel]; !ok {
		panic("Unknown rpc type")
	}

	convert2Opt(&opts, s.opt.beforeStart, micro.BeforeStart, true)
	convert2Opt(&opts, s.opt.beforeStop, micro.BeforeStop, false)
	convert2Opt(&opts, s.opt.afterStart, micro.AfterStart, false)
	convert2Opt(&opts, s.opt.afterStop, micro.AfterStop, false)
	opts = append(
		opts,
		micro.Server(cmd.DefaultServers[s.opt.RpcModel]()),
		micro.Client(cmd.DefaultClients[s.opt.RpcModel](
			client.Wrap(hystrix.NewClientWrapper(s.opt.EnableBreaker))),
		),
	)

	s.Service = micro.NewService(opts...)
	// TODO close config ,Later need to open
	_ = s.Options().Config.Close()
	return nil
}

// default init server
func (s *Service) initServer() (err error) {
	options := []server.Option{
		server.Name(s.opt.ServerName),
		server.Version(s.opt.ServerVersion),
		server.Address(s.opt.Address),
		server.Metadata(s.opt.metadata),
		server.Registry(s.register),
		server.RegisterTTL(time.Duration(s.opt.RegisterTTL) * time.Second),
		server.RegisterInterval(time.Duration(s.opt.RegisterInterval) * time.Second),
		server.WrapHandler(tracer.NewHandlerWrapper(s.opt.tracer)),
		server.WrapHandler(metrics.NewHandlerWrapper(s.opt.metrics)),
	}
	//server custom wrapper
	if len(s.opt.serverWrapper) > 0 {
		options = append(options, s.opt.serverWrapper...)
	}
	err = s.Server().Init(options...)

	if s.opt.serviceFunc != nil {
		s.opt.serviceFunc(s)
	}

	return
}

// default init client
func (s *Service) initClient() error {
	options := []client.Option{
		client.PoolSize(s.opt.PoolSize),
		client.PoolTTL(s.opt.PoolTTL),
		client.Retries(s.opt.Retries),
		client.RequestTimeout(s.opt.RequestTimeout),
		client.DialTimeout(s.opt.DialTimeout),
		client.Registry(s.register),
		client.WrapCall(tracer.NewCallWrapper(s.opt.tracer, s.Server().String)),
		client.Selector(selector.NewSelector(selector.Registry(s.register))),
	}
	//client custom wrapper
	if len(s.opt.clientWrapper) > 0 {
		options = append(options, s.opt.clientWrapper...)
	}
	err := s.Client().Init(options...)
	if err != nil {
		return err
	}
	return nil
}

func convert2Opt(opts *[]micro.Option, hooks []func(),
	microOption func(fn func() error) micro.Option, withoutRecovery bool) {
	if len(hooks) > 0 {
		for _, hook := range hooks {
			h := hook
			if withoutRecovery {
				*opts = append(*opts, microOption(func() error {
					h()
					return nil
				}))
			} else {
				*opts = append(*opts, microOption(wrapRecover(h)))
			}
		}
	}
}

func wrapRecover(f func()) func() error {
	return func() error {
		defer func() {
			if r := recover(); r != nil {
				xlog.Fatalln(r)
			}
		}()
		f()
		return nil
	}
}
