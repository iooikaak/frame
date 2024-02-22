package gins

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/iooikaak/frame/gin"
)

// Instance gins实例
var Instance *Server
var signalChan chan os.Signal

func init() {
	var err error
	Instance, err = New()
	if err != nil {
		err = errors.New("Gin Server 创建失败：" + err.Error())
		panic(err)
	}

	signalChan = make(chan os.Signal)
}

// Run 启动GinServer
func Run(conf *Config) {
	if conf.BroadcastIP == "" {
		conf.BroadcastIP = conf.IP
	}
	if conf.BroadcastPort <= 0 {
		conf.BroadcastPort = conf.Port
	}
	// 初始化服务配置
	Instance.Init(conf)

	//使用docker stop 命令去关闭Container时，该命令会发送SIGTERM 命令到Container主进程，让主进程处理该信号，关闭Container，如果在10s内，未关闭容器，Docker Damon会发送SIGKILL 信号将Container关闭。
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		Instance.Start()
	}()

	<-signalChan
	signal.Stop(signalChan)

	Instance.Stop()
}

// Stop 停止GinServer
func Stop() {
	signalChan <- syscall.SIGINT
}

// AddInit 添加安全初始化函数
func AddInit(initFuncs ...InitFunc) {
	Instance.AddInit(initFuncs...)
}

// AddTemplate 添加模板
func AddTemplate(name string, content string) {
	if Instance.templ == nil {
		Instance.templ = template.Must(template.New(name).Parse(content))
	} else {
		Instance.templ = template.Must(Instance.templ.New(name).Parse(content))
	}
}

// RouterGroup 路由对象
func RouterGroup(relativePath string, handlers ...HandlerFunc) *Router {
	return Instance.Router.Group(relativePath, handlers...)
}

// Use 全局中间件
func Use(middlewares ...HandlerFunc) {
	Instance.Middleware.Use(middlewares...)
}

// POST 请求
func POST(relativePath string, handlers ...HandlerFunc) *Router {
	return Instance.Router.POST(relativePath, handlers...)
}

// GET 请求
func GET(relativePath string, handlers ...HandlerFunc) *Router {
	return Instance.Router.GET(relativePath, handlers...)
}

// StaticFile 静态文件服务
func StaticFile(relativePath, filepath string) *Router {
	return Instance.Router.StaticFile(relativePath, filepath)
}

// Static 静态文件目录服务
func Static(relativePath, root string) *Router {
	return Instance.Router.Static(relativePath, root)
}

// StaticFS 静态资源服务
func StaticFS(relativePath string, fs http.FileSystem) *Router {
	return Instance.Router.StaticFS(relativePath, fs)
}

// On404 404自定义处理
func On404(handler HandlerFunc) {
	Instance.on404 = handler
}

// On500 500自定义处理
func On500(handler HandlerFunc) {
	Instance.on500 = handler
}

// Engine gin.Engine对象
// 本框架对Router进行了二次封装，调用此对象时请注意路由冲突问题
func Engine() *gin.Engine {
	return Instance.engine
}
