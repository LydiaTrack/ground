package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserCreateResponse struct {
	ID          bson.ObjectId `bson:"_id"`
	Username    string        `json:"username"`
	PersonInfo  `json:"personInfo"`
	CreatedDate time.Time `json:"createdDate,omitempty"`
	Version     int       `json:"version,omitempty"`
}
