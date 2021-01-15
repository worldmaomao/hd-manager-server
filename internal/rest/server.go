package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"worldmaomao/harddisk/internal/config"
	"worldmaomao/harddisk/internal/constant"
	"worldmaomao/harddisk/internal/rest/middlewares"
)

type server struct {
	container di.Container
}

func NewServer(container di.Container) *server {
	return &server{
		container: container,
	}
}

func (server *server) Start() {
	config := server.container.Get(constant.Configuration).(*config.Configuration)
	address := fmt.Sprintf("%s:%d", config.Service.Host, config.Service.Port)
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))
	// 不需要登录
	noAuthApiRoute := r.Group("/api/v1")
	loadPingRouter(noAuthApiRoute)
	loadAuthRouter(noAuthApiRoute, server.container)
	// 需要登录
	authApiRoute := r.Group("/api/v1")
	authApiRoute.Use(func(context *gin.Context) {
		middlewares.RequireAuthenticated(context, server.container)
	})
	loadDiskRouter(authApiRoute, server.container)
	loadDiskItemRouter(authApiRoute, server.container)
	r.Run(address)
}
