package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"net/http"
	"strconv"
	"worldmaomao/harddisk/internal/service"
	"worldmaomao/harddisk/internal/vo"
)

func getParameter(c *gin.Context, key string) string {
	var (
		result string
	)
	result, _ = c.GetPostForm(key)
	if len(result) == 0 {
		result, _ = c.GetQuery(key)
	}
	return result
}

func search(c *gin.Context, container di.Container) error {
	var (
		keyword        string
		isSearch       string
		page           int
		pageSize       int
		nextPage       int
		prePage        int
		searchCount    int
		searchItemList []vo.DiskItemVo
		itemService    = service.NewDiskItemService(container)
	)
	keyword = getParameter(c, "keyword")
	isSearch = getParameter(c, "isSearch")
	pageStr := getParameter(c, "page")
	pageSizeStr := getParameter(c, "pageSize")
	page, _ = strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSize, _ = strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 30
	}
	prePage = page - 1
	if prePage < 1 {
		prePage = 1
	}
	nextPage = page + 1
	if isSearch == "true" {
		searchCount, _ = itemService.QueryItemCount("", keyword)
		searchItemList, _ = itemService.QueryItem("", keyword, page, pageSize)
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":          "硬盘文件搜索",
		"isSearch":       isSearch,
		"keyword":        keyword,
		"page":           page,
		"prePage":        prePage,
		"nextPage":       nextPage,
		"pageSize":       pageSize,
		"searchItemList": searchItemList,
		"searchCount":    searchCount,
	})
	return nil
}
