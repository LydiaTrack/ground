package session

import "gopkg.in/mgo.v2/bson"

// InfoModel is a struct that contains the session information and maps to the userId
type InfoModel struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	UserId       bson.ObjectId `json:"userId" bson:"userId"`
	ExpireTime   int64         `json:"expireTime" bson:"expireTime"`
	RefreshToken string        `json:"refreshToken" bson:"refreshToken"`
}
