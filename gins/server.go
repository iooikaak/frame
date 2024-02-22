package gins

import (
	"context"
	"html/template"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	"github.com/iooikaak/frame/eureka"

	"github.com/iooikaak/frame/tracer"
	"github.com/iooikaak/frame/xlog"

	"github.com/iooikaak/frame/gin"
)

// InitFunc 安全初始化函数
type InitFunc func()

// Server 服务器对象
type Server struct {
	engine       *gin.Engine // gin Engine
	Middleware   *Middleware // 全局中间件
	Router       *Router     // gin Router 封装
	Tracer       *tracer.Tracer
	templ        *template.Template // 模板资源
	initFuncList []InitFunc         // 安全初始化函数列表

	server     *http.Server // http服务器
	rootCtx    context.Context
	rootCancel context.CancelFunc

	on404 HandlerFunc
	on500 HandlerFunc

	Config *Config

	Client *eureka.Client
}

// New 创建新的GinServer实例
func New() (gs *Server, err error) {
	gs = &Server{
		engine: gin.New(),
	}

	// 封闭自定义 Middleware ，全局
	gs.Middleware = &Middleware{engine: gs.engine}

	// 封装自定义 Router
	gs.Router = &Router{RouterGroup: &gs.engine.RouterGroup}

	gs.rootCtx, gs.rootCancel = context.WithCancel(context.Background())

	return
}

// Init 初始化
func (gs *Server) Init(conf *Config) {
	if conf.Name == "" {
		panic("name启动参数不能为空")
	}

	if conf.Version == "" {
		panic("version启动参数不能为空")
	}

	if conf.IP == "" {
		conf.IP = "0.0.0.0"
	}

	if conf.Port <= 0 {
		panic("port启动参数不能为空")
	}

	if conf.BroadcastIP == "" {
		conf.BroadcastIP = conf.IP
	}

	if conf.BroadcastPort <= 0 {
		conf.BroadcastPort = conf.Port
	}

	if conf.Timeout <= 0 {
		panic("timeout启动参数不能为空")
	}
	// 调试模式
	if conf.Debug {
		// 测试模式下 gin 的log看的清楚一点
		gs.engine.Use(gin.Logger())
	} else {
		if !conf.DisableAccessLog {
			gs.Middleware.Use(logger())
		}
		gin.SetMode(gin.ReleaseMode)
	}

	// 性能监测
	if conf.Pprof {
		pprofGroup := gs.engine.Group("/debug/pprof")

		pprofGroup.GET("/cmdline", func(ctx *gin.Context) {
			pprof.Cmdline(ctx.Writer, ctx.Request)
		})

		pprofGroup.GET("/profile", func(ctx *gin.Context) {
			pprof.Profile(ctx.Writer, ctx.Request)
		})

		pprofGroup.GET("/symbol", func(ctx *gin.Context) {
			pprof.Symbol(ctx.Writer, ctx.Request)
		})

		pprofGroup.GET("/trace", func(ctx *gin.Context) {
			pprof.Trace(ctx.Writer, ctx.Request)
		})

		pprofGroup.GET("/", func(ctx *gin.Context) {
			pprof.Index(ctx.Writer, ctx.Request)
		})

		pprofGroup.GET("/block", func(ctx *gin.Context) {
			pprof.Index(ctx.Writer, ctx.Request)
		})

		pprofGroup.GET("/goroutine", func(ctx *gin.Context) {
			pprof.Index(ctx.Writer, ctx.Request)
		})

		pprofGroup.GET("/heap", func(ctx *gin.Context) {
			pprof.Index(ctx.Writer, ctx.Request)
		})

		pprofGroup.GET("/mutex", func(ctx *gin.Context) {
			pprof.Index(ctx.Writer, ctx.Request)
		})

		pprofGroup.GET("/threadcreate", func(ctx *gin.Context) {
			pprof.Index(ctx.Writer, ctx.Request)
		})
	}

	// 设置 http server
	gs.server = &http.Server{
		Addr:    conf.Addr(),
		Handler: gs.engine,
		//ReadTimeout: conf.ReadTimeout,
		//WriteTimeout: conf.WriteTimeout,
	}

	gs.Config = conf

	//添加gin链路
	//if opentracing.IsGlobalTracerRegistered() && gs.Tracer != nil {
	//gs.engine.Use(tracer.NewGinJaegerTrace())
	//}

	// 加载核心中间件
	gs.engine.Use(core())

	// 加载全局中间件
	gs.Middleware.init()

	// 初始化路由
	gs.Router.init()

	// 加载HTML模板
	// gin.Engine 在创建时，模板尚未初始完毕，需要在这里再进行设置
	// 因 Start 只会调用一次，在 Stop 后应用会直接退出，忽略线程不安全的警告
	if gs.templ != nil {
		gs.engine.SetHTMLTemplate(gs.templ)
	}

	// 加载安全初始化函数
	for _, fn := range gs.initFuncList {
		fn()
	}
}

// AddInit 添加安全初始化函数
func (gs *Server) AddInit(initFunc ...InitFunc) {
	if len(initFunc) > 0 {
		gs.initFuncList = append(gs.initFuncList, initFunc...)
	}
}

func (gs *Server) AddTracer(t *tracer.Tracer) {
	gs.Tracer = t
	gs.Middleware.Use(NewGinJaegerTrace(t.Instance()))
}

// Start 启动服务
func (gs *Server) Start() {
	xlog.Infof("[%s %s]服务运行在：%s", gs.Config.Name, gs.Config.Version, gs.Config.BroadcastAddr())

	err := gs.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		xlog.Errorf("%s", err.Error())
	}
}

// Stop 停止服务
func (gs *Server) Stop() {
	xlog.Infof("正在停止[%s %s]服务...", gs.Config.Name, gs.Config.Version)

	stopCtx, stopCancel := context.WithTimeout(gs.rootCtx, time.Duration(gs.Config.Timeout)*time.Second)

	// FIXME: 不关闭的话，优雅退出时，会导致有连接挂起，总是需要超时退出,目前关闭 keep-alive 状态并未成功  2019年8月
	// FIXME: 但是关闭 KeepAlive 性能受影响太多,决定暂时放开 2020年1月
	gs.server.SetKeepAlivesEnabled(true)

	err := gs.server.Shutdown(stopCtx)

	if err != nil {
		xlog.Warnf("[%s %s]服务停止出错：%s", gs.Config.Name, gs.Config.Version, err)
	} else {
		xlog.Infof("[%s %s]服务已停止", gs.Config.Name, gs.Config.Version)
	}

	// 关闭 tracer reporter
	gs.rootCancel()
	if gs.Tracer != nil {
		if err := gs.Tracer.Close(); err != nil {
			xlog.Error(err)
		}
	}
	stopCancel()

	// 延时2秒退出，让超时任务 504 响应完成
	time.Sleep(2 * time.Second)
}

func (gs *Server) Engine() *gin.Engine {
	return gs.engine
}
func (gs *Server) InitEureka(conf *Config) {
	//create eureka client
	if !strings.Contains(conf.Name, "gins") {
		conf.Name = conf.Name + ".gins"
	}
	gs.Client = eureka.NewClient(&eureka.Config{
		DefaultZone:                    conf.CenterAddr,
		App:                            conf.Name,
		Port:                           conf.Port,
		RenewalIntervalInSecs:          conf.RenewalIntervalInSecs,
		DurationInSecs:                 conf.DurationInSecs,
		RegistryFetchIntervalSeconds:   conf.RegistryFetchIntervalSeconds,
		RollDiscoveriesIntervalSeconds: conf.RollDiscoveriesIntervalSeconds,
		DataCenterName:                 conf.CenterName,
	})
	conf.HeathUrl = gs.Client.Config.Instance.HealthCheckURL
	conf.StatusUrl = gs.Client.Config.Instance.StatusPageURL
}
