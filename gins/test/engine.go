package test

import (
	"github.com/iooikaak/frame/gin"
	"github.com/iooikaak/frame/gins"
)

// Engine 服务对象
var Engine *gin.Engine

// Init 初始化gins Engine
func Init(conf *gins.Config) {
	gins.Instance.Init(conf)

	Engine = gins.Engine()
}
