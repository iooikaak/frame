package job

import (
	"context"
	"fmt"
	"github.com/iooikaak/frame/net/ip"
	"github.com/matthew188/go-xxl-job"
	"github.com/matthew188/go-xxl-job/handler"
	"github.com/matthew188/go-xxl-job/logger"
)

//执行器
type XxlJob struct {
	isInit bool
}

func InitExecutor(xxlConfig *XxlJobConfig) (xxlClient *XxlJob, err error) {
	xxlClient = new(XxlJob)
	if len(xxlConfig.Addresses) == 0 {
		err = fmt.Errorf("xxl-job: addresses cannot be empty")
		return
	}

	if len(xxlConfig.AppName) == 0 {
		err = fmt.Errorf("xxl-job: app_name cannot be empty")
		return
	}

	if xxlConfig.Port <= 0 {
		p, e := ip.GetFreePort("")
		if e != nil {
			panic(e.Error())
		}
		xxlConfig.Port = p
	}
	//注册执行器
	xxl.InitExecutor(xxlConfig.Addresses, xxlConfig.AccessToken, xxlConfig.AppName, xxlConfig.Port)
	return
}
func InitExecutorFromJson(config string) (xxlClient *XxlJob, err error) {
	xxlClient = new(XxlJob)
	xxlConfig, err := newConfig(config)
	if err != nil {
		return
	}
	if xxlConfig.Port <= 0 {
		p, e := ip.GetFreePort("")
		if e != nil {
			panic(e.Error())
		}
		xxlConfig.Port = p
	}
	//注册执行器
	xxl.InitExecutor(xxlConfig.Addresses, xxlConfig.AccessToken, xxlConfig.AppName, xxlConfig.Port)

	return
}
func (x *XxlJob) RegisterJob(jobName string, job handler.JobHandlerFunc) {
	//注册job
	xxl.RegisterJob(jobName, job)
	x.isInit = true
}
func (x *XxlJob) RunServer() {
	if !x.isInit {
		panic("尚未注册")
	}
	//启动server
	go xxl.RunServer()
	fmt.Println(`xxl_job执行器已经启动。`)
}
func GetParam(ctx context.Context, key string) (val string, has bool) {
	val, has = xxl.GetParam(ctx, key)
	return
}
func Info(ctx context.Context, args ...interface{}) {
	logger.Info(ctx, args)
}
