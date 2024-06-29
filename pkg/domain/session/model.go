package session

import "go.mongodb.org/mongo-driver/bson/primitive"

// InfoModel is a struct that contains the session information and maps to the userId
type InfoModel struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserId       primitive.ObjectID `json:"userId" bson:"userId"`
	ExpireTime   int64              `json:"expireTime" bson:"expireTime"`
	RefreshToken string             `json:"refreshToken" bson:"refreshToken"`
}
