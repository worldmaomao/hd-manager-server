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
)

type DiskService struct {
	diskDao   *dao.DiskDao
	container di.Container
}

func NewDiskService(container di.Container) *DiskService {
	return &DiskService{
		container: container,
		diskDao:   container.Get(constant.DiskDao).(*dao.DiskDao),
	}
}

func (service *DiskService) GetById(id string) (*model.DiskModel, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("the id for disk is illegal")
	}
	return service.diskDao.GetById(bson.ObjectIdHex(id))
}

func (service *DiskService) GetByIds(ids []string) ([]model.DiskModel, error) {
	if ids == nil || len(ids) == 0 {
		return nil, fmt.Errorf("ids is empty")
	}
	bsonObjectIds := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		if !bson.IsObjectIdHex(id) {
			return nil, fmt.Errorf("the format of is illegal")
		}
		bsonObjectIds[i] = bson.ObjectIdHex(id)
	}
	return service.diskDao.GetByIds(bsonObjectIds)
}

func (service *DiskService) GetMapByIds(ids []string) (map[string]model.DiskModel, error) {
	models, err := service.GetByIds(ids)
	if err != nil {
		return nil, err
	}
	modelMap := make(map[string]model.DiskModel)
	for _, model := range models {
		modelMap[model.Id.Hex()] = model
	}
	return modelMap, nil
}

func (service *DiskService) DeleteDisk(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("illeal disk id[%s]", id)
	}
	diskItemService := NewDiskItemService(service.container)
	if count, _ := diskItemService.QueryItemCount(id, ""); count > 0 {
		return fmt.Errorf("该硬盘里还有文件存在，不能删除")
	}
	diskObjectId := bson.ObjectIdHex(id)
	if disk, _ := service.diskDao.GetById(diskObjectId); disk == nil {
		return fmt.Errorf("该硬盘已经删除")
	}
	return service.diskDao.Delete([]bson.ObjectId{diskObjectId})
}

func (service *DiskService) QueryDisk(page int, pageSize int, name string) ([]model.DiskModel, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return service.diskDao.QueryList(model.DiskModel{
		Name: name,
	}, (page-1)*pageSize, pageSize)
}

func (service *DiskService) QueryDiskCount(key string) (int, error) {
	return service.diskDao.QueryCount(model.DiskModel{
		Name: key,
	})
}

func (service *DiskService) UpdateDisk(newDisk model.DiskModel) error {
	var (
		err error
	)
	v, _ := utils.NewValidator()
	err = v.Struct(newDisk)
	if err != nil {
		return err
	}
	existDisk, _ := service.diskDao.GetById(newDisk.Id)
	if existDisk == nil {
		return fmt.Errorf("硬盘没有找到")
	}
	existDiskByName, _ := service.diskDao.GetByName(newDisk.Name)
	if existDiskByName != nil && existDiskByName.Id.Hex() != newDisk.Id.Hex() {
		return fmt.Errorf("该硬盘名称已存在，请换一个")
	}
	existDisk.Name = newDisk.Name
	existDisk.Description = newDisk.Description
	existDisk.Modified = time.Now().Unix()
	return service.diskDao.Update(*existDisk)
}

func (service *DiskService) AddDisk(disk model.DiskModel) error {
	var (
		err error
	)
	v, _ := utils.NewValidator()
	err = v.Struct(disk)
	if err != nil {
		return err
	}
	existDisk, err := service.diskDao.GetByName(disk.Name)
	if err == nil && existDisk != nil {
		return errors.NewDuplicationError(fmt.Sprintf("磁盘[%s]已经存在", disk.Name))
	}
	disk.Created = time.Now().Unix()
	disk.Modified = time.Now().Unix()
	return service.diskDao.Add(disk)
}
