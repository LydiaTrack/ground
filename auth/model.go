package auth

import "gopkg.in/mgo.v2/bson"

type Permission struct {
	Domain string `json:"domain"`
	Action string `json:"action"`
}

type AuthContext struct {
	Permissions []Permission  `json:"permissions"`
	UserId      bson.ObjectId `json:"userId"`
}
