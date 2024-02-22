package bootstrap

//
//import (
//	"context"
//	"testing"
//	"time"
//
//	"github.com/iooikaak/frame/tracer"
//
//	demoPB "github.com/iooikaak/proto/demoAdmin"
//	"github.com/iooikaak/frame/gins"
//	"github.com/iooikaak/frame/micro"
//	"github.com/iooikaak/frame/xlog"
//)
//
//func TestRun(t *testing.T) {
//	// 启动框架
//	app := New(Config{
//		Log: &xlog.Config{
//			Service:   "mf.demoadmin",
//			Host:      "",
//			Caller:    false,
//			Stdout:    true,
//			File:      nil,
//			NsqConfig: nil,
//		},
//		Tracer: &tracer.TracingConfig{
//			TracingType:  "jaeger",
//			SamplerType:  "const",
//			SamplerParam: "1",
//			SenderType:   "udp",
//			Endpoint:     "10.180.18.20:6831",
//		},
//		GinServer: &gins.Config{
//			Version:                        "0.0.1",
//			Host:                           "localhost",
//			IP:                             "0.0.0.0",
//			BroadcastIP:                    "",
//			Port:                           4041,
//			BroadcastPort:                  0,
//			Timeout:                        5,
//			Debug:                          true,
//			Pprof:                          true,
//			ReadTimeout:                    time.Second,
//			WriteTimeout:                   5 * time.Second,
//			DisableAccessLog:               true,
//			CenterAddr:                     []string{"127.0.0.1:8761"},
//			CenterName:                     "discovery",
//			RenewalIntervalInSecs:          5,
//			RegistryFetchIntervalSeconds:   2,
//			RollDiscoveriesIntervalSeconds: 5,
//			DurationInSecs:                 10,
//		},
//		MicroServer: &micro.Options{
//			PoolSize:                       100,
//			Retries:                        1,
//			PoolTTL:                        10,
//			RequestTimeout:                 5 * time.Second,
//			DialTimeout:                    5 * time.Second,
//			RpcRegister:                    "eureka",
//			RpcModel:                       "grpc",
//			CenterAddr:                     []string{"127.0.0.1:8761", "127.0.0.2:8762", "127.0.0.1:8763"},
//			CenterName:                     "discovery",
//			RegisterTTL:                    30,
//			RegisterInterval:               2,
//			RenewalIntervalInSecs:          2,
//			RollDiscoveriesIntervalSeconds: 5,
//			DisableTracer:                  false,
//			EnableBreaker:                  false,
//			ServerVersion:                  "1.0.1",
//			Address:                        ":10087",
//			ServerOnly:                     true,
//		},
//		ServiceName: "mf.demoadmin",
//	})
//
//	// 注册服务, grpc可选
//	app.Init(
//		//http process
//		HTTPService(func(r *gins.Server) {
//			//http.Init(conf.Conf, r)
//		}),
//		MicroService(func(s *micro.Service) {
//			if err := demoPB.RegisterDemoAdminHandler(s.Server(), NewService()); err != nil {
//				panic(err)
//			}
//		}),
//	)
//
//	if err := app.Run(); err != nil {
//		panic(err)
//	}
//}
//
//type Service struct {
//}
//
//func NewService() *Service {
//	return &Service{}
//}
//
//func (s *Service) Operate(ctx context.Context, input *demoPB.OperateInput, output *demoPB.OperateOutput) (err error) {
//	output.Message = "成功了"
//	output.Code = demoPB.ErrorCode_Success
//	return
//}
