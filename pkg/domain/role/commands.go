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
	Name        string            `json:"name"`
	Info        string            `json:"info,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Permissions []auth.Permission `json:"permissions"`
}

type DeleteRoleCommand struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}
