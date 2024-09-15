package role

import (
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateRoleCommand struct {
	Name        string            `json:"name"`
	Tags        []string          `json:"tags,omitempty"`
	Info        string            `json:"info,omitempty"`
	Permissions []auth.Permission `json:"permissions"`
}

type UpdateRoleCommand struct {
	Name        string            `json:"name" bson:"name"`
	Info        string            `json:"info,omitempty" bson:"info,omitempty"`
	Tags        []string          `json:"tags,omitempty" bson:"tags,omitempty"`
	Permissions []auth.Permission `json:"permissions" bson:"permissions"`
}

type DeleteRoleCommand struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}
