package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserCreateResponse struct {
	ID          bson.ObjectId `json:"_id"`
	Username    string        `json:"username"`
	PersonInfo  `json:"person_info"`
	CreatedDate time.Time `json:"created_date,omitempty"`
	Version     int       `json:"version,omitempty"`
}
