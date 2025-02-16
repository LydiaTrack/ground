package role

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LydiaTrack/ground/pkg/auth"
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

type Option func(*Model) error

func NewRole(opts ...Option) (*Model, error) {
	u := &Model{
		ID:          primitive.NewObjectID(),
		CreatedDate: time.Now(),
		Version:     1,
	}

	for _, opt := range opts {
		if err := opt(u); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func WithName(name string) Option {
	return func(r *Model) error {
		r.Name = name
		return nil
	}
}

func WithPermissions(permissions []auth.Permission) Option {
	return func(r *Model) error {
		r.Permissions = permissions
		return nil
	}
}

func WithTags(tags []string) Option {
	return func(r *Model) error {
		r.Tags = tags
		return nil
	}
}

func WithInfo(info string) Option {
	return func(r *Model) error {
		r.Info = info
		return nil
	}
}

func (r Model) Validate() error {

	if len(r.Name) == 0 {
		return errors.New("name is required")
	}

	return nil
}

func (r Model) HasPermissions(permissions []auth.Permission) bool {
	for _, permission := range permissions {
		if !r.HasPermission(permission) {
			return false
		}
	}

	return true
}

func (r Model) HasPermission(permission auth.Permission) bool {
	for _, p := range r.Permissions {
		if p == permission {
			return true
		}
	}

	return false
}
