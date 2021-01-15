package vo

import "worldmaomao/harddisk/internal/dao/model"

type DiskItemVo struct {
	model.DiskItemModel
	DiskName string `json:"diskName"`
}
