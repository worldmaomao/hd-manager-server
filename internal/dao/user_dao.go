package dao

import (
	"gopkg.in/mgo.v2/bson"
	"sync"
	"worldmaomao/harddisk/internal/dao/model"
	"worldmaomao/harddisk/internal/database/interfaces"
)

var (
	userDaoInstant *UserDao
	userDaoOnce    sync.Once
)

type UserDao struct {
	dbClient       interfaces.DbClient
	collectionName string
}

func NewUserDao(dbClient interfaces.DbClient) *UserDao {
	userDaoOnce.Do(func() {
		userDaoInstant = &UserDao{
			dbClient:       dbClient,
			collectionName: "users",
		}
	})
	return userDaoInstant
}

func (dao *UserDao) GetByUsername(username string) (*model.User, error) {
	session := dao.dbClient.GetSession()
	var user model.User
	err := session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Find(bson.M{"username": username}).One(&user)
	return &user, err
}

func (dao *UserDao) Add(user model.User) error {
	session := dao.dbClient.GetSession()
	defer session.Close()
	return session.DB(dao.dbClient.GetConfig().DatabaseName).C(dao.collectionName).Insert(user)
}
