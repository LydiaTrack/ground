package role

import (
	"errors"
	"lydia-track-base/internal/domain/auth"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	ID          bson.ObjectId     `json:"id" bson:"_id"`
	Name        string            `json:"name" bson:"name"`
  Permissions []auth.Permission `json:"permissions" bson:"permissions"`
	Tags        []string          `json:"tags,omitempty" bson:"tags,omitempty"`
	RoleInfo    string            `json:"roleInfo,omitempty" bson:"roleInfo,omitempty"`
	CreatedDate time.Time         `json:"createdDate" bson:"createdDate"`
	Version     int               `json:"version" bson:"version"`
}

func NewRole(id string, name string, tags []string, roleInfo string, createdDate time.Time, version int) Model {
	return Model{
		ID:          bson.ObjectIdHex(id),
		Name:        name,
		Permissions: []auth.Permission{},
		Tags:        tags,
		RoleInfo:    roleInfo,
		CreatedDate: createdDate,
		Version:     version,
	}
}

func (r Model) Validate() error {

	if len(r.Name) == 0 {
		return errors.New("rolename is required")
	}

	return nil
}
