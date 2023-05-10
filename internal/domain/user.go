package domain

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID         bson.ObjectId `bson:"_id"`
	Password   string        `bson:"password"`
	PersonInfo `bson:"person_info"`
}

func NewUser(id string, password string, personInfo PersonInfo) User {
	return User{
		ID:         bson.ObjectIdHex(id),
		Password:   password,
		PersonInfo: personInfo,
	}
}
