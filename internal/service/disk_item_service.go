package service

import (
	"fmt"
	"github.com/sarulabs/di"
	"gopkg.in/mgo.v2/bson"
	"time"
	"worldmaomao/harddisk/internal/constant"
	"worldmaomao/harddisk/internal/dao"
	"worldmaomao/harddisk/internal/dao/model"
	"worldmaomao/harddisk/internal/errors"
	"worldmaomao/harddisk/internal/utils"
	"worldmaomao/harddisk/internal/vo"
)

type DiskItemService struct {
	container   di.Container
	diskItemDao *dao.DiskItemDao
}

func NewDiskItemService(container di.Container) *DiskItemService {
	return &DiskItemService{
		container:   container,
		diskItemDao: container.Get(constant.DiskItemDao).(*dao.DiskItemDao),
	}
}

func (service *DiskItemService) CheckExist(diskId string, diskItemName string) (bool, error) {
	if len(diskId) == 0 || len(diskItemName) == 0 {
		return false, errors.NewParameterError("参数错误")
	}
	model, err := service.diskItemDao.GetByName(diskId, diskItemName)
	return model != nil, err
}

func (service *DiskItemService) AddItem(item model.DiskItemModel) error {
	var (
		err error
	)
	v, _ := utils.NewValidator()
	if err = v.Struct(item); err != nil {
		return err
	}
	diskService := NewDiskService(service.container)
	disk, _ := diskService.GetById(item.DiskId)
	if disk == nil {
		return fmt.Errorf("fail to find disk by id[%s]", item.DiskId)
	}
	existDiskItem, _ := service.diskItemDao.GetByName(item.DiskId, item.FileName)
	if existDiskItem != nil {
		return fmt.Errorf("该文件[%s]已经入库", item.FileName)
	}
	item.Modified = time.Now().Unix()
	item.Created = time.Now().Unix()
	return service.diskItemDao.Add(item)
}

func (service *DiskItemService) UpdateItem(item model.DiskItemModel) error {
	var (
		err error
	)
	v, _ := utils.NewValidator()
	if err = v.Struct(item); err != nil {
		return err
	}
	diskService := NewDiskService(service.container)
	disk, err := diskService.GetById(item.DiskId)
	if err != nil {
		return err
	}
	if disk == nil {
		return fmt.Errorf("fail to find disk by id[%s]", item.DiskId)
	}
	existDiskItem, _ := service.diskItemDao.GetByName(item.DiskId, item.FileName)
	if existDiskItem != nil && existDiskItem.Id.Hex() != item.Id.Hex() {
		return fmt.Errorf("该文件的文件名[%s]已经被占用", item.FileName)
	}
	item.Modified = time.Now().Unix()
	return service.diskItemDao.Update(item)
}

func (service *DiskItemService) DeleteItem(itemId string) error {
	if !bson.IsObjectIdHex(itemId) {
		return fmt.Errorf("illegal disk item id")
	}
	itemObjectId := bson.ObjectIdHex(itemId)
	existItem, _ := service.diskItemDao.GetById(itemObjectId)
	if existItem == nil {
		return fmt.Errorf("该文件已经被删除")
	}
	return service.diskItemDao.Delete([]bson.ObjectId{itemObjectId})
}

func (service *DiskItemService) QueryItem(diskId string, name string, page int, pageSize int) ([]vo.DiskItemVo, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	diskModelList, err := service.diskItemDao.QueryList(model.DiskItemModel{
		DiskId:   diskId,
		FileName: name,
	}, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	diskIdList := make([]string, 0)
	diskIdMap := make(map[string]bool)
	for _, model := range diskModelList {
		if _, ok := diskIdMap[model.DiskId]; !ok {
			diskIdMap[model.DiskId] = true
			diskIdList = append(diskIdList, model.DiskId)
		}
	}
	diskService := NewDiskService(service.container)
	disModelMap, err := diskService.GetMapByIds(diskIdList)
	if err != nil {
		return nil, err
	}

	diskVoList := make([]vo.DiskItemVo, len(diskModelList))
	for i, model := range diskModelList {
		var diskName string
		if _, ok := disModelMap[model.DiskId]; ok {
			diskName = disModelMap[model.DiskId].Name
		} else {
			diskName = "unknown"
		}
		var vo = vo.DiskItemVo{
			DiskItemModel: model,
			DiskName:      diskName,
		}
		diskVoList[i] = vo
	}
	return diskVoList, nil
}

func (service *DiskItemService) QueryItemCount(diskId string, name string) (int, error) {
	return service.diskItemDao.QueryCount(model.DiskItemModel{
		DiskId:   diskId,
		FileName: name})
}
