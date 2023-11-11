package session

import (
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/auth"
)

// InfoModel is a struct that contains the session information and maps to the userId
type InfoModel struct {
	ID           bson.ObjectId `bson:"_id"`
	UserId       bson.ObjectId `bson:"userId"`
	ExpireTime   int64         `bson:"expireTime"`
	RefreshToken string        `bson:"refreshToken"`
}

// UserSession is a struct that contains the user specific session information.
// It is calculating when user sends any request to the server.
type UserSession struct {
	Permissions []auth.Permission
}
