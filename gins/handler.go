package gins

import "github.com/iooikaak/frame/gin"

// HandlerFunc 路由方法
type HandlerFunc func(*Context)

// newGinHandler 重新封装
func newGinHandler(handlers ...HandlerFunc) []gin.HandlerFunc {
	l := len(handlers)
	if l <= 0 {
		return nil
	}

	ginHandlers := make([]gin.HandlerFunc, 0, l)

	for i := 0; i < l; i++ {
		handler := handlers[i]

		ginHandlers = append(ginHandlers, func(ginCtx *gin.Context) {
			// recovery 全局中间件的使用，保证了 *Context 的存在
			val, _ := ginCtx.Get("*Context")
			ctx := val.(*Context)

			// 执行原 HandlerFunc
			handler(ctx)
		})
	}

	return ginHandlers
}
