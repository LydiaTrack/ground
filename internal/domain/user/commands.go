package user

import (
	"github.com/Lydia/lydia-base/internal/domain/auth"
	"gopkg.in/mgo.v2/bson"
)

type CreateUserCommand struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	PersonInfo `json:"personInfo"`
}

type UpdateUserCommand struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	PersonInfo `json:"personInfo"`
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

var CreatePermission = auth.Permission{
	Domain: "user",
	Action: "CREATE",
}

var UpdatePermission = auth.Permission{
	Domain: "user",
	Action: "UPDATE",
}

var DeletePermission = auth.Permission{
	Domain: "user",
	Action: "DELETE",
}

var ReadPermission = auth.Permission{
	Domain: "user",
	Action: "READ",
}
