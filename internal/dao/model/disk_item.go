package model

import (
	"gopkg.in/mgo.v2/bson"
)

type DiskItemModel struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	DiskId      string        `bson:"diskId" json:"diskId" validate:"required" label:"硬盘Id"`
	FileName    string        `bson:"fileName" json:"fileName" validate:"required,max=1024,min=1" label:"文件名称"`
	PicPath     string        `bson:"picPath,omitempty" json:"picPath"`
	FileType    string        `bson:"fileType" json:"fileType" validate:"required,max=20,min=1" label:"文件类型"`
	Description string        `bson:"description,omitempty" json:"description" label:"描述"`
	Created     int64         `bson:"created" json:"created"`
	Modified    int64         `bson:"modified" json:"modified"`
}
