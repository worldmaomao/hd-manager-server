package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
	"worldmaomao/harddisk/internal/dao/model"
	"worldmaomao/harddisk/internal/errors"
	"worldmaomao/harddisk/internal/service"
	"worldmaomao/harddisk/internal/vo"
)

func loadDiskItemRouter(group *gin.RouterGroup, container di.Container) {
	group.POST("/disk-item", NewHandlerWrapper(container).handle(itemAdd))
	group.PUT("/disk-item", NewHandlerWrapper(container).handle(itemUpdate))
	group.DELETE("/disk-item", NewHandlerWrapper(container).handle(itemDelete))
	group.GET("/disk-item", NewHandlerWrapper(container).handle(itemQuery))
	group.GET("/disk-item/exist", NewHandlerWrapper(container).handle(checkItemExist))
}

func checkItemExist(c *gin.Context, container di.Container) error {
	var (
		diskId   string
		fileName string
		service  = service.NewDiskItemService(container)
	)
	diskId = c.DefaultQuery("diskId", "")
	fileName = c.DefaultQuery("fileName", "")
	if len(diskId) == 0 || len(fileName) == 0 {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	exist, err := service.CheckExist(diskId, fileName)
	if err != nil {
		return errors.NewRestError(http.StatusOK, 0, err.Error())
	}
	c.JSON(http.StatusOK, vo.Response{
		Code:    0,
		Message: "",
		Data:    exist,
	})
	return nil
}

func itemAdd(c *gin.Context, container di.Container) error {
	var (
		diskItemArray []model.DiskItemModel
		resultVos     []vo.ImportDiskItemResult
		service       = service.NewDiskItemService(container)
	)
	if err := c.ShouldBindJSON(&diskItemArray); err != nil {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	for _, item := range diskItemArray {
		if err := service.AddItem(item); err != nil {
			resultVos = append(resultVos, vo.ImportDiskItemResult{
				DiskItemName: item.FileName,
				Error:        err.Error(),
			})
		}
	}
	c.JSON(http.StatusCreated, vo.Response{
		Code:    0,
		Message: "",
		Data:    resultVos,
	})
	return nil
}

func itemUpdate(c *gin.Context, container di.Container) error {
	var (
		err     error
		json    model.DiskItemModel
		service = service.NewDiskItemService(container)
	)
	if err = c.ShouldBindJSON(&json); err != nil {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	if err = service.UpdateItem(json); err != nil {
		log.Printf("fail to update disk, error:%s", err.Error())
		return errors.NewRestError(http.StatusInternalServerError, -1, err.Error())
	}
	c.JSON(http.StatusAccepted, vo.Response{
		Code:    0,
		Message: "更新成功",
	})
	return nil
}

func itemDelete(c *gin.Context, container di.Container) error {
	var (
		diskItemIds string
		resultVos   []vo.DeleteDiskItemResult
		service     = service.NewDiskItemService(container)
	)
	diskItemIds = c.DefaultQuery("id", "")
	diskItemIdArray := strings.Split(diskItemIds, ",")
	if diskItemIdArray == nil || len(diskItemIdArray) == 0 {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	for _, diskItemId := range diskItemIdArray {
		if !bson.IsObjectIdHex(diskItemId) {
			return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
		}
	}
	for _, diskItemId := range diskItemIdArray {
		if err := service.DeleteItem(diskItemId); err != nil {
			log.Printf("fail to delete disk, error:%s", err.Error())
			resultVos = append(resultVos, vo.DeleteDiskItemResult{
				DiskItemId: diskItemId,
				Error:      err.Error(),
			})
		}
	}
	c.JSON(http.StatusAccepted, vo.Response{
		Code:    0,
		Message: "",
		Data:    resultVos,
	})
	return nil
}

func itemQuery(c *gin.Context, container di.Container) error {
	var (
		err         error
		page        string
		pageInt     int
		pageSize    string
		pageSizeInt int
		diskId      string
		fileName    string
		service     = service.NewDiskItemService(container)
	)
	page = c.DefaultQuery("page", "1")
	pageSize = c.DefaultQuery("pageSize", "10")
	diskId = c.DefaultQuery("diskId", "")
	fileName = c.DefaultQuery("fileName", "")
	if pageInt, err = strconv.Atoi(page); err != nil {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	if pageSizeInt, err = strconv.Atoi(pageSize); err != nil {
		return errors.NewRestError(http.StatusBadRequest, -1, "参数错误")
	}
	diskList, _ := service.QueryItem(diskId, fileName, pageInt, pageSizeInt)
	count, _ := service.QueryItemCount(diskId, fileName)
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
