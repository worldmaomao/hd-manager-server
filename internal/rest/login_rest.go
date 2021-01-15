package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"net/http"
	"worldmaomao/harddisk/internal/errors"
	"worldmaomao/harddisk/internal/service"
	"worldmaomao/harddisk/internal/vo"
)

func loadNoAuthRouter(group *gin.RouterGroup, container di.Container) {
	group.POST("/login", NewHandlerWrapper(container).handle(login))
}

func login(c *gin.Context, container di.Container) error {
	var loginVo vo.LoginVo
	if err := c.ShouldBindJSON(&loginVo); err != nil {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	userService := service.NewUserService(container)
	jwt, err := userService.Login(loginVo.Username, loginVo.Password, "web")
	if err != nil {
		return errors.NewRestError(http.StatusCreated, -1, err.Error())
	}
	c.JSON(http.StatusCreated, vo.Response{
		Code:    0,
		Message: "",
		Data:    jwt,
	})
	return nil
}
