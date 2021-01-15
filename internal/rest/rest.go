package rest

import (
	"github.com/gin-gonic/gin"
)

func loadPingRouter(group *gin.RouterGroup) {
	group.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.String(200, "pong")
}
