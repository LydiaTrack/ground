package commands

import (
	"gopkg.in/mgo.v2/bson"
)

type CreateRoleCommand struct {
	RoleName    string 		  `json:"rolename"`
	Tags        []string      `json:"tags,omitempty"`
	RoleInfo    string        `bson:"role_info,omitempty"`
}

type UpdateRoleCommand struct {
	RoleName string `json:"rolename"`
}

type DeleteRoleCommand struct {
	ID bson.ObjectId `json:"_id"`
}
