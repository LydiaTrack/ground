package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserCreateResponse struct {
	ID          bson.ObjectId `bson:"_id"`
	Username    string        `bson:"username"`
	PersonInfo  `bson:"person_info"`
	CreatedDate time.Time `bson:"created_date,omitempty"`
	Version     int       `bson:"version,omitempty"`
}
