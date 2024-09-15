package role

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/LydiaTrack/lydia-base/pkg/auth"
)

type Model struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Permissions []auth.Permission  `json:"permissions" bson:"permissions"`
	Tags        []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	Info        string             `json:"info,omitempty" bson:"info,omitempty"`
	CreatedDate time.Time          `json:"createdDate" bson:"createdDate"`
	Version     int                `json:"version" bson:"version"`
}

func NewRole(id string, name string, permissions []auth.Permission, tags []string, info string, createdDate time.Time, version int) Model {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Model{}
	}
	return Model{
		ID:          objID,
		Name:        name,
		Tags:        tags,
		Permissions: permissions,
		Info:        info,
		CreatedDate: createdDate,
		Version:     version,
	}
}

func (r Model) Validate() error {

	if len(r.Name) == 0 {
		return errors.New("name is required")
	}

	return nil
}
