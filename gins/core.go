package gins

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/iooikaak/frame/apiconstant"
	"github.com/iooikaak/frame/gin"
	"github.com/iooikaak/frame/xlog"
)

var ctxPool sync.Pool

func init() {
	ctxPool.New = func() interface{} {
		return &Context{}
	}
}

// core panic恢复，初始化Context
func core() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ctxPool.New().(*Context)
		ctx.reset(ginCtx)
		ctx.Set("*Context", ctx)
		defer ctxPool.Put(ctx)
		defer func() {
			// 异常捕获处理
			if e := recover(); e != nil {
				stack := fmt.Sprintf("System Panic: %v", e)
				for i := 1; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					} else {
						stack += "\n"
					}

					stack += fmt.Sprintf("%v:%v", file, line)
				}

				var rawReq []byte
				if ginCtx.Request != nil {
					rawReq, _ = httputil.DumpRequest(ginCtx.Request, false)
				}
				pl := fmt.Sprintf("%s ERROR:http call panic: %s\n%v\n%s\n", time.Now().Format("2006-01-02 15:04:05"), string(rawReq), e, stack)
				_, _ = fmt.Fprint(os.Stderr, pl)
				l := xlog.WithContext(ginCtx)
				l.Logger.ExitFunc = func(code int) {}
				l.Fatalln(pl)

				// 500 处理
				if Instance.on500 != nil {
					ctx.setPanic(stack)
					Instance.on500(ctx)
				} else {
					// 默认异常响应
					if ctx.IsAPI() {
						ctx.API.SetError(NewAPIErrorWithLog("系统异常", stack))
						ctx.API.Render()
						ctx.Abort()
					} else {
						ctx.AbortWithStatus(http.StatusInternalServerError)
					}
				}
			}

		}()

		ctx.Next()
		if ctx.IsAborted() {
			return
		}

		status := ctx.Writer.Status()

		// 路由匹配到的情况下，status默认 200
		if status != http.StatusNotFound && ctx.IsAPI() {
			if ctx.API.result.Code == apiconstant.RESPONSE_UNKNOW {
				ctx.API.result.Msg = "API空响应"
			}

			ctx.API.Render()
		}

		if status == http.StatusNotFound {
			// 404 处理
			if Instance.on404 != nil {
				Instance.on404(ctx)
			} else {
				ctx.AbortWithStatus(404)
			}
		}
	}
}
