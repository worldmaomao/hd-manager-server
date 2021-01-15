package dao

import (
	"gopkg.in/mgo.v2/bson"
	"sync"
	"worldmaomao/harddisk/internal/dao/model"
	"worldmaomao/harddisk/internal/database/interfaces"
)

var (
	diskDaoInstant *DiskDao
	diskDaoOnce    sync.Once
)

type DiskDao struct {
	dbClient       interfaces.DbClient
	collectionName string
}

func NewDiskDao(dbClient interfaces.DbClient) *DiskDao {
	diskDaoOnce.Do(func() {
		diskDaoInstant = &DiskDao{
			dbClient:       dbClient,
			collectionName: "disks",
		}
	})
	return diskDaoInstant
}

func (dao *DiskDao) GetById(bsonId bson.ObjectId) (*model.DiskModel, error) {
	session := dao.dbClient.GetSession()
	defer session.Close()
	var disk model.DiskModel
	err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(bson.M{"_id": bsonId}).One(&disk)
	if err != nil {
		return nil, err
	}
	return &disk, nil
}

func (dao *DiskDao) GetByIds(bsonIds []bson.ObjectId) ([]model.DiskModel, error) {
	session := dao.dbClient.GetSession()
	defer session.Close()
	var disks []model.DiskModel
	err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(bson.M{"_id": bson.M{"$in": bsonIds}}).All(&disks)
	if err != nil {
		return nil, err
	}
	return disks, nil
}

func (dao *DiskDao) GetByName(name string) (*model.DiskModel, error) {
	session := dao.dbClient.GetSession()
	defer session.Close()
	var disk model.DiskModel
	err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(bson.M{"name": name}).One(&disk)
	if err != nil {
		return nil, err
	}
	return &disk, nil
}

func (dao *DiskDao) Add(disk model.DiskModel) error {
	session := dao.dbClient.GetSession()
	defer session.Close()
	return session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Insert(disk)
}

func (dao *DiskDao) Update(disk model.DiskModel) error {
	session := dao.dbClient.GetSession()
	defer session.Close()
	return session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Update(bson.M{"_id": disk.Id}, disk)
}

func (dao *DiskDao) Delete(bsonIds []bson.ObjectId) error {
	session := dao.dbClient.GetSession()
	defer session.Close()
	_, err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).RemoveAll(bson.M{"_id": bson.M{"$in": bsonIds}})
	return err
}

func (dao *DiskDao) createQuery(disk model.DiskModel) interface{} {
	bsons := make([]bson.M, 0)
	if len(disk.Name) > 0 {
		bsons = append(bsons, bson.M{"name": bson.RegEx{Pattern: disk.Name, Options: "i"}})
	}
	var query interface{}
	if len(bsons) > 0 {
		query = bson.M{"$and": bsons}
	}
	return query
}

func (dao *DiskDao) QueryCount(disk model.DiskModel) (int, error) {
	session := dao.dbClient.GetSession()
	defer session.Close()
	return session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(dao.createQuery(disk)).Count()
}

func (dao *DiskDao) QueryList(disk model.DiskModel, start int, limit int) ([]model.DiskModel, error) {
	session := dao.dbClient.GetSession()
	defer session.Close()
	var disks []model.DiskModel
	err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(dao.createQuery(disk)).Skip(start).Limit(limit).Sort("-created").All(&disks)
	return disks, err
}
