package commands

import (
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/auth"
	"lydia-track-base/internal/domain/user"
)

type CreateUserCommand struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	user.PersonInfo `json:"person_info"`
}

type UpdateUserCommand struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	user.PersonInfo `json:"person_info"`
}

type DeleteUserCommand struct {
	ID bson.ObjectId `json:"_id"`
}

type AddRoleToUserCommand struct {
	UserID bson.ObjectId `json:"user_id"`
	RoleID bson.ObjectId `json:"role_id"`
}

type RemoveRoleFromUserCommand struct {
	UserID bson.ObjectId `json:"user_id"`
	RoleID bson.ObjectId `json:"role_id"`
}

var CreatePermission = auth.Permission{
	Domain: "user",
	Action: "CREATE",
}
