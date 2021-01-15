package dao

import (
	"gopkg.in/mgo.v2/bson"
	"sync"
	"worldmaomao/harddisk/internal/dao/model"
	"worldmaomao/harddisk/internal/database/interfaces"
)

var (
	diskItemDaoInstant *DiskItemDao
	diskItemDaoOnce    sync.Once
)

type DiskItemDao struct {
	dbClient       interfaces.DbClient
	collectionName string
}

func NewDiskItemDao(dbClient interfaces.DbClient) *DiskItemDao {
	diskItemDaoOnce.Do(func() {
		diskItemDaoInstant = &DiskItemDao{
			dbClient:       dbClient,
			collectionName: "disk_items",
		}
	})
	return diskItemDaoInstant
}

func (dao *DiskItemDao) GetById(bsonId bson.ObjectId) (*model.DiskItemModel, error) {
	session := dao.dbClient.GetSession()
	defer session.Close()
	var diskItem model.DiskItemModel
	err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(bson.M{"_id": bsonId}).One(&diskItem)
	if err != nil {
		return nil, err
	}
	return &diskItem, nil
}

func (dao *DiskItemDao) GetByName(diskId string, name string) (*model.DiskItemModel, error) {
	session := dao.dbClient.GetSession()
	defer session.Close()
	var disk model.DiskItemModel
	err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(bson.M{"fileName": name, "diskId": diskId}).One(&disk)
	if err != nil {
		return nil, err
	}
	return &disk, nil
}

func (dao *DiskItemDao) Add(diskItem model.DiskItemModel) error {
	session := dao.dbClient.GetSession()
	defer session.Close()
	return session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Insert(diskItem)
}

func (dao *DiskItemDao) Update(diskItem model.DiskItemModel) error {
	session := dao.dbClient.GetSession()
	defer session.Close()
	return session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Update(bson.M{"_id": diskItem.Id}, diskItem)
}

func (dao *DiskItemDao) Delete(bsonIds []bson.ObjectId) error {
	session := dao.dbClient.GetSession()
	defer session.Close()
	_, err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).RemoveAll(bson.M{"_id": bson.M{"$in": bsonIds}})
	return err
}

func (dao *DiskItemDao) createQuery(diskItem model.DiskItemModel) interface{} {
	bsons := make([]bson.M, 0)
	if len(diskItem.FileName) > 0 {
		bsons = append(bsons, bson.M{"fileName": bson.RegEx{Pattern: diskItem.FileName, Options: "i"}})
	}
	if len(diskItem.DiskId) > 0 {
		bsons = append(bsons, bson.M{"diskId": diskItem.DiskId})

	}
	var query interface{}
	if len(bsons) > 0 {
		query = bson.M{"$and": bsons}
	}
	return query
}

func (dao *DiskItemDao) QueryCount(diskItem model.DiskItemModel) (int, error) {
	session := dao.dbClient.GetSession()
	defer session.Close()
	return session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(dao.createQuery(diskItem)).Count()
}

func (dao *DiskItemDao) QueryList(diskItem model.DiskItemModel, start int, limit int) ([]model.DiskItemModel, error) {
	session := dao.dbClient.GetSession()
	defer session.Close()
	var diskItems []model.DiskItemModel
	err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(dao.createQuery(diskItem)).Skip(start).Limit(limit).Sort("-created").All(&diskItems)
	return diskItems, err
}
