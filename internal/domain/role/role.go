package role

import (
	"errors"
	"lydia-track-base/internal/domain/auth"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	ID          bson.ObjectId     `bson:"_id"`
	Name        string            `bson:"name"`
	Permissions []auth.Permission `bson:"permissions"`
	Tags        []string          `bson:"tags,omitempty"`
	RoleInfo    string            `bson:"role_info,omitempty"`
	CreatedDate time.Time         `bson:"created_date,omitempty"`
	Version     int               `bson:"version,omitempty"`
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
