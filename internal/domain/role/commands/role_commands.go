package commands

import (
	"lydia-track-base/internal/domain/auth"

	"gopkg.in/mgo.v2/bson"
)

type CreateRoleCommand struct {
	Name string   `json:"name"`
	Tags []string `json:"tags,omitempty"`
	Info string   `json:"roleInfo,omitempty"`
}

type UpdateRoleCommand struct {
	Name string `json:"name"`
}

type DeleteRoleCommand struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

var CreatePermission = auth.Permission{
	Domain: "role",
	Action: "CREATE",
}

var UpdatePermission = auth.Permission{
	Domain: "role",
	Action: "UPDATE",
}

var DeletePermission = auth.Permission{
	Domain: "role",
	Action: "DELETE",
}

var ReadPermission = auth.Permission{
	Domain: "role",
	Action: "READ",
}
