package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"log"
	"net/http"
	"strconv"
	"worldmaomao/harddisk/internal/dao/model"
	"worldmaomao/harddisk/internal/errors"
	"worldmaomao/harddisk/internal/service"
	"worldmaomao/harddisk/internal/vo"
)

func loadDiskRouter(group *gin.RouterGroup, container di.Container) {
	group.POST("/disk", NewHandlerWrapper(container).handle(diskAdd))
	group.PUT("/disk", NewHandlerWrapper(container).handle(diskUpdate))
	group.DELETE("/disk", NewHandlerWrapper(container).handle(diskDelete))
	group.GET("/disk", NewHandlerWrapper(container).handle(diskQuery))
}

func diskAdd(c *gin.Context, container di.Container) error {
	var (
		err     error
		json    model.DiskModel
		service = service.NewDiskService(container)
	)
	if err := c.ShouldBindJSON(&json); err != nil {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	if err = service.AddDisk(json); err != nil {
		log.Printf("fail to add disk, error:%s", err.Error())
		return errors.NewRestError(http.StatusCreated, -1, err.Error())
	}
	c.JSON(http.StatusCreated, vo.Response{
		Code:    0,
		Message: "添加成功",
	})
	return nil
}

func diskUpdate(c *gin.Context, container di.Container) error {
	var (
		err     error
		json    model.DiskModel
		service = service.NewDiskService(container)
	)
	if err = c.ShouldBindJSON(&json); err != nil {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	if err = service.UpdateDisk(json); err != nil {
		log.Printf("fail to update disk, error:%s", err.Error())
		return errors.NewRestError(http.StatusInternalServerError, -1, err.Error())
	}
	c.JSON(http.StatusAccepted, vo.Response{
		Code:    0,
		Message: "更新成功",
	})
	return nil
}

func diskDelete(c *gin.Context, container di.Container) error {
	var (
		diskId  string
		err     error
		service = service.NewDiskService(container)
	)
	diskId = c.DefaultQuery("id", "")
	if err = service.DeleteDisk(diskId); err != nil {
		log.Printf("fail to delete disk, error:%s", err.Error())
		return errors.NewRestError(http.StatusAccepted, -1, err.Error())
	}
	c.JSON(http.StatusAccepted, vo.Response{
		Code:    0,
		Message: "删除成功",
	})
	return nil
}

func diskQuery(c *gin.Context, container di.Container) error {
	var (
		err         error
		page        string
		pageInt     int
		pageSize    string
		pageSizeInt int
		keyword     string
		service     = service.NewDiskService(container)
	)
	page = c.DefaultQuery("page", "1")
	pageSize = c.DefaultQuery("pageSize", "10")
	keyword = c.DefaultQuery("keyword", "")
	if pageInt, err = strconv.Atoi(page); err != nil {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	if pageSizeInt, err = strconv.Atoi(pageSize); err != nil {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	diskList, _ := service.QueryDisk(pageInt, pageSizeInt, keyword)
	count, _ := service.QueryDiskCount(keyword)
	c.JSON(http.StatusOK, vo.Response{
		Code:    0,
		Message: "",
		Data: vo.Page{
			Page:  pageInt,
			Size:  pageSizeInt,
			Total: count,
			List:  diskList,
		},
	})
	return nil
}
