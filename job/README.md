# Job

xxl-job的golang版执行器实现，我们提供了基于 `matthew188/go-xxl-job-client` 的封装

## 定时任务注册

文件位于 ``job/xxl_job`` 文件中,下面会引入一个例子:

```golang
# demo task
package main

import (
	"github.com/iooikaak/frame/job"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	xxlConfig1 := &job.XxlJobConfig{
		Addresses: []string{"http://10.180.18.81:8080/xxl-job-admin"},
		AppName:   "local-job-executor",
	}

	xxlClient, err := job.InitExecutor(xxlConfig1)
	if err != nil {
		panic(err)
	}
	xxlClient.RegisterJob("ShDemoJob", ShDemoJob)
	xxlClient.RegisterJob("ShDemoJob2", ShDemoJob2)
	xxlClient.RegisterJob("ShDemoJob3", ShDemoJob3)

	xxlClient.RunServer()
}
func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-signals
}

func ShDemoJob(ctx context.Context) error {
	job.Info(ctx, "golang job1 run success >>>>>>>>>>>>>>")
	job.Info(ctx, "the input param:", )
	v, e := job.GetParam(ctx, "name")
	job.Info(ctx, v, " has,", e)
	return nil
}

func ShDemoJob2(ctx context.Context) error {
	job.Info(ctx, "golang job2 run success >>>>>>>>>>>>>>")
	job.Info(ctx, "the input param:", )
	v, e := job.GetParam(ctx, "name")
	job.Info(ctx, v, " has,", e)
	return nil
}
func ShDemoJob3(ctx context.Context) error {
	job.Info(ctx, "golang job3 run success >>>>>>>>>>>>>>")
	job.Info(ctx, "the input param:", )
	v, e := job.GetParam(ctx, "name")
	job.Info(ctx, v, " has,", e)
	return nil
}



```

上面例子中, ``job.InitExecutor``或``job.InitExecutorFromJson`` 会初始化执行器，app_name为执行器名称,初始化完成后可至 ``Addresses``
地址对应的管理页面：http://job.juqitech.com/ 添加相应名称的执行器

> ``job.Info(ctx, "some logs")`` ：该方法打印日志至后台

> ``job.GetParam(ctx, "name")`` ：该方法获取``key``为``name``的执行参数,在后台添加job时多参数逗号隔开，形如： ``key1=val1,key2=val2``

> ``xxlClient.RegisterJob`` 注册job,可注册多个

> ``JobHandlerFunc``格式必须为：``func(ctx context.Context) error ``

最后启动服务``xxlClient.RunServer()``，默认以goroutine方式与框架同运行(在框架main函数中引入即可)，若要作为单独项目运行，则需注意防止main函数结束：

```
signals := make(chan os.Signal, 1)
signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
<-signals
```
