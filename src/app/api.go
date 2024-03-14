package app

import (
	"context"
	"errors"
	"gin-mvc/src/core"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	engine   = gin.Default()
	router   *core.CRouter
	ApiRoot  = "/api/v1"
	Accounts = gin.Accounts{
		"admin": "admin",
		"guest": "guest",
	}
)

// DemoServer 优雅停机服务
type DemoServer struct {
	HttpSrv *http.Server
	TimeOut time.Duration
}

func (agent *DemoServer) Start() {
	log.Printf("Listening and serving HTTP on %s\n", agent.HttpSrv.Addr)
	err := agent.HttpSrv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen %s\n", err)
	}
}

func (agent *DemoServer) Stop(ctx context.Context) error {
	return agent.HttpSrv.Shutdown(ctx)
}

func (agent *DemoServer) Timeout() time.Duration {
	return agent.TimeOut
}

func Run(addr string) {
	// 设置router
	router = &core.CRouter{IRouter: engine.Group(ApiRoot)}
	// register middleware
	router.RegisterMiddleware(gin.BasicAuth(Accounts))
	// register routers
	// 注册方式: Route
	router.RegisterRoute(core.CRoute{Path: "/mvc1", Controller: &Controller{}})
	// 注册方式: 自动扫描
	router.RegisterController("/mvc2", &Controller{})
	// start engine
	server := &DemoServer{
		HttpSrv: &http.Server{
			Addr:         addr,
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		TimeOut: 15 * time.Second,
	}
	core.StartSecureServer(server)
}
