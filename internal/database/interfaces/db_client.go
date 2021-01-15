package interfaces

import (
	"gopkg.in/mgo.v2"
	"worldmaomao/harddisk/internal/config"
)

type DbClient interface {
	GetSession() *mgo.Session
	GetConfig() config.Database
}
