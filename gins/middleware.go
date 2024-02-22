package gins

import (
	"github.com/iooikaak/frame/gin"
)

// Middleware 全局中间件
type Middleware struct {
	engine   *gin.Engine
	handlers []HandlerFunc
}

// init 注册全局中间件到gin中
func (m *Middleware) init() {
	if len(m.handlers) <= 0 {
		return
	}

	ginHandlers := newGinHandler(m.handlers...)
	m.engine.Use(ginHandlers...)
}

// Use 添加全局中间件
func (m *Middleware) Use(middlewares ...HandlerFunc) {
	m.handlers = append(m.handlers, middlewares...)
}
