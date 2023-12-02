package commands

import (
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain/user"
)

type CreateUserCommand struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	user.PersonInfo `json:"personInfo"`
}

type UpdateUserCommand struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	user.PersonInfo `json:"personInfo"`
}

type DeleteUserCommand struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

type AddRoleToUserCommand struct {
	UserID bson.ObjectId `json:"userId"`
	RoleID bson.ObjectId `json:"roleId"`
}

type RemoveRoleFromUserCommand struct {
	UserID bson.ObjectId `json:"userId"`
	RoleID bson.ObjectId `json:"roleId"`
}
