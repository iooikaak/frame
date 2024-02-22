package gins

import (
	"strings"
)

// Action Action风格路由中间件
// 在RouterGroup下使用中间件时，必须 路由路径 匹配时才会触发，全局使用时无此问题
func Action(prefixPath ...string) HandlerFunc {
	return func(ctx *Context) {
		query := ctx.Request.URL.Query()

		action := query.Get("action")
		if action == "" {
			// 非Action风格，退回
			ctx.Next()
			return
		}

		actions := strings.Split(action, ".")

		// 增加跳转地址前缀
		urlPath := "/" + strings.Join(actions, "/") + "/"
		if len(prefixPath) == 1 {
			urlPath = prefixPath[0] + urlPath
		}

		// 更新目标
		query.Del("action")
		ctx.Request.URL.Path = urlPath
		ctx.Request.URL.RawQuery = query.Encode()
		Instance.engine.HandleContext(ctx.Context)

		// 命中action风格处理，中止后续行为
		ctx.Abort()
	}
}
