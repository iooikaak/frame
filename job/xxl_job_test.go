package job

import (
	"context"
	"testing"
)

func TestXxlJob(t *testing.T) {
	xxlConfig1 := &XxlJobConfig{
		Addresses: []string{"http://10.180.18.81:8080/xxl-job-admin"},
		AppName:   "local-job-executor",
	}

	//1.初始化执行器
	xxlClient, err := InitExecutor(xxlConfig1)
	if err != nil {
		t.Error(err)
		return
	}
	//2.注册job
	xxlClient.RegisterJob("ShDemoJob0", ShDemoJob)
	xxlClient.RegisterJob("ShDemoJob1", ShDemoJob2)

	//3.运行server
	xxlClient.RunServer()
}
func TestXxlJobJson(t *testing.T) {
	xxlConfig1 := `{
  	"adds": ["http://10.180.18.81:8080/xxl-job-admin"],
  	"access_token": "",
  	"app_name": "local-job-executor"
	}`

	xxlClient, err := InitExecutorFromJson(xxlConfig1)
	if err != nil {
		t.Error(err)
		return
	}
	xxlClient.RegisterJob("ShDemoJob3", ShDemoJob)
	xxlClient.RegisterJob("ShDemoJob4", ShDemoJob2)

	xxlClient.RunServer()
}
func ShDemoJob(ctx context.Context) error {
	Info(ctx, "golang job1 run success >>>>>>>>>>>>>>")
	Info(ctx, "the input param:")
	v, e := GetParam(ctx, "name")
	Info(ctx, v, " has,", e)
	return nil
}

func ShDemoJob2(ctx context.Context) error {
	Info(ctx, "golang job2 run success >>>>>>>>>>>>>>")
	Info(ctx, "the input param:")
	v, e := GetParam(ctx, "name")
	Info(ctx, v, " has,", e)
	return nil
}
