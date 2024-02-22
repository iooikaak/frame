package gins

import (
	"bytes"
	"context"
	"fmt"
	"github.com/iooikaak/frame/gin"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContextDiffContext(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("gins.Context 不匹配原生上下文, 请检查异常: %s", err)
		}
	}()
	// 这里检测下 gins 上下文是否适配 context.Context(), 同时翻阅过源代码,http原生的context默认是一个空的上下文
	// http 上下文源码引用参见官方源码: https://golang.org/src/net/http/request.go?m=text
	// 所以,只要覆盖了 gins.Context 的默认 Value 方法,就适配了原生上下文
	var _ context.Context = &Context{}
	t.Log("上下文兼容性检查成功,可以正常使用")
}

func TestContextGetValue(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	ginsCtx := &Context{Context: c}
	ginsCtx.Set("AxA", 1)
	fmt.Println(ginsCtx.Get("AxA"))
}
