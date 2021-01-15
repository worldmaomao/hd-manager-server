package middlewares

import (
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"net/http"
	"worldmaomao/harddisk/internal/vo"
)

func NewRecovery() gin.HandlerFunc {
	return nice.Recovery(recoveryHandler)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, vo.Response{
		Code:    -1,
		Message: "服务器内部错误",
	})
}
