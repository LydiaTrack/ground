package commands

import (
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
