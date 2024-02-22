package xlog

import (
	"context"
	"testing"

	"github.com/iooikaak/frame/metadata"

	"github.com/pkg/errors"
)

type contextKey string

// init stdout log
func init() {
	err := Init(&Config{
		Service: "service_test",
		// Host:      "",
		Stdout: true,
		Caller: true,
		// Dir:       "",
		// NsqConfig: nil,
	})
	if err != nil {
		panic(err)
	}
}

// init dir file log
// func init() {
//	err := Init(&Config{
//		Service: "service_test",
//		Dir:     "/data/log/test.log"})
//	if err != nil {
//		panic(err)
//	}
// }

// init NSQ log
// func init() {
//	err := Init(&Config{
//		Service: "service_test",
//		NsqConfig: &NsqConfig{
//			addr:  "http://10.180.18.20:4161",
//			topic: "servicelog",
//		}})
//	if err != nil {
//		panic(err)
//	}
// }

func TestFileHook(t *testing.T) {
	_namespace = "dev"
	log := NewLogrus(&Config{
		Service: "service_test",
		File: &DirConfig{
			Dir:        "/tmp",
			Type:       "trace",
			MaxSize:    500,
			MaxBackups: 1,
			MaxAge:     1,
		}})
	log.WithField("param", "fff").Errorf("hehe %s", "hahaha")
}

func TestInfof(t *testing.T) {
	Infof("111 %s%d", "infof", 123)
}

func TestWarnf(t *testing.T) {
	Warnf("111 %s", "warnf")
}

func TestErrorf(t *testing.T) {
	Errorf("111 %s", "errorf")
}

func TestInfo(t *testing.T) {
	Info("info", 123)
}

func TestWarn(t *testing.T) {
	Warn("warn")
}

func TestErrorWithFields(t *testing.T) {
	_ = Init(&Config{
		Service: "test",
		Stdout:  true,
		Caller:  true,
	})

	vs := make(map[string]interface{})
	vs["name"] = "testname"
	vs["age"] = 10
	vs["occupation"] = "golang"

	// ErrorWithFields(vs, "error")
	InfoWithFields(vs, "error")
}

func TestInfoWithFields(t *testing.T) {
	vs := make(map[string]interface{})
	vs["name"] = "testname"
	vs["age"] = 10
	vs["occupation"] = "golang"
	InfoWithFields(vs, "info", 123)
}

func TestLogWithFields(t *testing.T) {
	ctx := context.Background()
	vs := make(map[string]interface{})
	vs["name"] = "testname"
	vs["age"] = 10
	vs["occupation"] = "golang"
	WithContext(ctx).WithFields(vs).Log(DebugLevel)
}

func TestWithEventAndErr(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextKey(metadata.HttpTraceId), "default")
	WithContext(ctx).WithError(errors.New("i'm here")).WithEvent("test").Info("hi error")
}

func TestWithErr(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextKey(metadata.HttpTraceId), "default")
	err := errors.New("error1")
	err = errors.Wrap(err, "error2_")
	WithContext(ctx).WithError(err).Info("have_error")
}

func TestWithError(t *testing.T) {
	err := errors.New("error1")
	err = errors.Wrap(err, "error2_")
	WithError(err).Info("check_error...")
}

func TestWithEvent(t *testing.T) {
	WithEvent("testEvent").Info("event_click")
}

func TestWarnWithFields(t *testing.T) {
	vs := make(map[string]interface{})
	vs["name"] = "testname"
	vs["age"] = 10
	vs["occupation"] = "golang"
	WarnWithFields(vs, "info", 123)
}

func TestErrorWithFieldsFmt(t *testing.T) {
	vs := make(map[string]interface{})
	vs["name"] = "testname"
	vs["age"] = 10
	vs["occupation"] = "golang"
	ErrorWithFieldsFmt(vs, "%s%d", "error", 1)
}

func TestInfoWithFieldsFmt(t *testing.T) {
	vs := make(map[string]interface{})
	vs["name"] = "testname"
	vs["age"] = 10
	vs["occupation"] = "golang"
	InfoWithFieldsFmt(vs, "%s%d", "info", 123)
}

func TestWarnWithFieldsFmt(t *testing.T) {
	vs := make(map[string]interface{})
	vs["name"] = "testname"
	vs["age"] = 10
	vs["occupation"] = "golang"
	WarnWithFields(vs, "info", 123)
}

func BenchmarkLogWithFields(b *testing.B) {
	vs := make(map[string]interface{})
	vs["name"] = "default"
	vs["age"] = 10
	vs["occupation"] = "golang"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InfoWithFields(vs, "logger benchmark")
	}
}

func BenchmarkLoggerWithFields(b *testing.B) {
	vs := make(map[string]interface{})
	vs["name"] = "default"
	vs["age"] = 10
	vs["occupation"] = "golang"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithFields(vs).Info("logger benchmark")
	}
}
