package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"time"
	"worldmaomao/harddisk/internal/config"
	"worldmaomao/harddisk/internal/database/interfaces"
)

type dbClientImpl struct {
	session        *mgo.Session
	databaseConfig config.Database
}

func (client *dbClientImpl) GetSession() *mgo.Session {
	return client.session.Copy()
}

func (client *dbClientImpl) GetConfig() config.Database {
	return client.databaseConfig
}

func NewDbClient(databaseConfig config.Database) (interfaces.DbClient, error) {
	var (
		mongoAddress    = fmt.Sprintf("%s:%d", databaseConfig.Host, databaseConfig.Port)
		mongoDBDialInfo *mgo.DialInfo
		session         *mgo.Session
		err             error
	)
	mongoDBDialInfo = &mgo.DialInfo{
		Addrs:    []string{mongoAddress},
		Timeout:  10 * time.Second,
		Database: databaseConfig.DatabaseName,
		Username: databaseConfig.Username,
		Password: databaseConfig.Password,
	}
	// 	mongoUrl := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", databaseConfig.Username, databaseConfig.Password, databaseConfig.Host, databaseConfig.Port, databaseConfig.DatabaseName)
	session, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, fmt.Errorf("connect db error:" + err.Error())
	}
	log.Println("database connected")
	return &dbClientImpl{
		session:        session,
		databaseConfig: databaseConfig,
	}, nil
}
