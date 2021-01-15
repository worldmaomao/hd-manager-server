package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"net/http"
	"worldmaomao/harddisk/internal/errors"
	"worldmaomao/harddisk/internal/vo"
)

type HandlerFunc func(c *gin.Context, di di.Container) error

type HandlerWrapper struct {
	di di.Container
}

func NewHandlerWrapper(di di.Container) *HandlerWrapper {
	return &HandlerWrapper{
		di: di,
	}
}

func (wrapper *HandlerWrapper) handle(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err error
		)
		err = handler(c, wrapper.di)
		if err != nil {
			switch err.(type) {
			case errors.RestError:
				err1 := err.(errors.RestError)
				c.JSON(err1.HttpResponseCode, vo.Response{
					Code:    err1.ErrorCode,
					Message: err1.ErrorMsg,
				})
			case error:
				c.JSON(http.StatusInternalServerError, vo.Response{
					Code:    -1,
					Message: err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, vo.Response{
					Code:    -1,
					Message: "未知错误",
				})
			}
		}
	}
}
