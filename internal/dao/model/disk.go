package model

import "gopkg.in/mgo.v2/bson"

type DiskModel struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name" json:"name" validate:"required,max=20,min=1" label:"硬盘名称"`
	Description string        `bson:"description" json:"description" validate:"max=255,min=0" label:"描述"`
	Created     int64         `bson:"created" json:"created"`
	Modified    int64         `bson:"modified" json:"modified"`
}
