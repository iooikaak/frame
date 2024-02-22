package bootstrap

import (
	"net/http"
	"time"

	"github.com/iooikaak/frame/gins"
)

type Info struct {
	Git   Git   `json:"git"`
	Build Build `json:"build"`
}

type Commit struct {
	Time time.Time `json:"time"`
	ID   string    `json:"id"`
}

type Git struct {
	Commit Commit `json:"commit"`
	Branch string `json:"branch"`
}

type Build struct {
	Version  string    `json:"version"`
	Artifact string    `json:"artifact"`
	Name     string    `json:"name"`
	Group    string    `json:"group"`
	Time     time.Time `json:"time"`
}

func (a *App) initGin() {
	// health
	a.httpServer.Router.GET("/actuator/health", func(ctx *gins.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	// info
	a.httpServer.Router.GET("/actuator/info", func(ctx *gins.Context) {
		ctx.JSON(Info{}, nil)
	})

	// http reg metric
	a.httpServer.AddMetric(a.metrics)

	// http reg tracing
	if a.tracer != nil {
		a.httpServer.AddTracer(a.tracer)
	}

	a.httpServer.InitEureka(a.conf.GinServer)

	// process
	if a.engineFunc != nil {
		a.engineFunc(a.httpServer)
	}

	a.httpServer.Init(a.conf.GinServer)

	a.beforeStart = append(a.beforeStart, func() {
		go func() {
			a.httpServer.Start()
		}()
	})

	a.beforeStart = append(a.beforeStart, func() {
		go func() {
			a.httpServer.Client.Start()
		}()
	})

	a.afterStop = append(a.afterStop, func() {
		a.httpServer.Stop()
	})
}
