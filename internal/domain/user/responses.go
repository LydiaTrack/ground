package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserCreateResponse struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Username    string        `json:"username" bson:"username"`
	PersonInfo  `json:"personInfo"`
	CreatedDate time.Time `json:"createdDate" bson:"createdDate"`
	Version     int       `json:"version" bson:"version"`
}
