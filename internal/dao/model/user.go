package model

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username string        `bson:"username" json:"username" validate:"required" label:"用户名"`
	Password string        `bson:"password" json:"password" validate:"required" label:"密码"`
	Roles    []string      `bson:"roles" json:"roles"`
}
