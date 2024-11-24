package auth

import "go.mongodb.org/mongo-driver/bson/primitive"

type Permission struct {
	Domain string `json:"domain"`
	Action string `json:"action"`
}

type PermissionContext struct {
	Permissions []Permission        `json:"permissions"`
	UserID      *primitive.ObjectID `json:"userID"`
}
