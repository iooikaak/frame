package xlog

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewLogrusInfo(t *testing.T) {
	log := NewLogrus(&Config{Service: "app_service"})
	log.Log(context.Background(), logrus.InfoLevel, 2020)
}

func TestNewLogrusInfof(t *testing.T) {
	log := NewLogrus(&Config{Service: "app_service"})

	log.Logf(context.Background(), logrus.InfoLevel, "%s %d", "hello!", 2020)
}

func TestNewLogrusInfofWithFields(t *testing.T) {
	log := NewLogrus(&Config{Service: "app_service"})
	vs := make(map[string]interface{})
	vs["name"] = "testname"
	vs["age"] = 10
	vs["occupation"] = "golang"
	log.LogWithFields(context.Background(), logrus.InfoLevel, vs)
}

func TestNewLogrusInfofWithFieldsFmt(t *testing.T) {
	log := NewLogrus(&Config{Service: "app_service"})
	vs := make(map[string]interface{})
	vs["name"] = "testname"
	vs["age"] = 10
	vs["occupation"] = "golang"
	log.LogWithFieldsFmt(context.Background(), logrus.InfoLevel, vs, "%s", "just for test")
}

//func TestNewLogrusLogDirInfo(t *testing.T) {
//	log := NewLogrus(&Config{
//		Dir:     "/data/log/test.log",
//		Service: "test_service",
//	})
//
//	log.Log(context.Background(), logrus.InfoLevel, "hello!", 2020)
//}
//
//func TestNewLogrusLogDir(t *testing.T) {
//	log := NewLogrus(&Config{
//		Dir:     "/data/log/test.log",
//		Service: "test_service",
//	})
//
//	log.Logf(context.Background(), logrus.InfoLevel, "%s %d", "hello!", 2020)
//}

func TestNewLogrusStdoutFmt(t *testing.T) {
	log := NewLogrus(&Config{Service: "app_service"})
	log.Logf(context.Background(), logrus.InfoLevel, "%s %d", "hello!", 2020)
}

//func TestNewLogrusLogDirFmt(t *testing.T) {
//	_namespace = "dev"
//	log := NewLogrus(&Config{Dir: "/tmp/xxx.log"})
//	for i := 0; i < 30; i++ {
//		log.Infof("%s %d", "hello 自证清白!", 2020)
//		time.Sleep(time.Second)
//	}
//}

func TestNewLogrusLogNsq(t *testing.T) {
	log := NewLogrus(&Config{
		Service: "test_service",
		NsqConfig: &NsqConfig{
			Addr:  "http://10.180.18.20:4161",
			Topic: "servicelog",
		}})

	log.Logf(context.Background(), logrus.InfoLevel, "%s %d %d", "hello!", 2020, time.Now().Unix())
}
