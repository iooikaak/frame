package example

import (
	"time"

	"github.com/iooikaak/frame/gins"
)

var Instance *gins.Server

// ExampleNew 实例化一个 gins 实例
func ExampleNew() {
	Instance, err := gins.New()
	if err != nil {
		panic(err)
	}
	Instance.Init(&gins.Config{
		Name:             "example",
		Version:          "1.0.0",
		Host:             "localhost",
		IP:               "127.0.0.1",
		Port:             4040,
		Timeout:          5,
		Debug:            false,
		Pprof:            false,
		ReadTimeout:      time.Second * 1,
		WriteTimeout:     time.Second * 1,
		DisableAccessLog: false,
	})
	// 启动, 暂时注释
	// Instance.Start()
}

func ExampleAddMiddleware() {
	Instance.Middleware.Use(func(c *gins.Context) {
		// todo something
		c.Next()
		// todo something
	})
}

func ExampleAddRouter() {
	Instance.Router.GET("/test", func(c *gins.Context) {
		c.API.SetDataKV("hello", "world")
	})
}

func ExampleRouterGroup() {
	g := Instance.Router.Group("/test")
	g.GET("/test", func(c *gins.Context) {
		c.API.SetDataKV("hello", "world")
	})
}

func ExampleRouterMiddlerware() {
	g := Instance.Router.Group("/test")
	g.Use(func(c *gins.Context) {
		// todo something
		c.Next()
		// todo something
	})
	g.GET("/test", func(c *gins.Context) {
		c.API.SetDataKV("hello", "world")
	})
}
