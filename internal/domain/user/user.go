package user

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain/role"
	"time"
)

type Option func(m Model) Model

type Model struct {
	ID          bson.ObjectId `bson:"_id"`
	Username    string        `bson:"username"`
	Password    string        `bson:"password"`
	PersonInfo  `bson:"person_info"`
	CreatedDate time.Time    `bson:"created_date,omitempty"`
	Version     int          `bson:"version,omitempty"`
	Roles       []role.Model `bson:"roles,omitempty"`
}

func NewUser(id string, username string, password string, personInfo PersonInfo, createdDate time.Time, version int, option ...Option) Model {
	m := Model{
		ID:          bson.ObjectIdHex(id),
		Username:    username,
		Password:    password,
		PersonInfo:  personInfo,
		CreatedDate: createdDate,
		Version:     version,
	}

	for _, opt := range option {
		m = opt(m)
	}

	return m
}

func (u Model) Validate() error {

	if len(u.Password) == 0 {
		return errors.New("password is required")
	}

	if len(u.Username) == 0 {
		return errors.New("username is required")
	}

	if err := u.PersonInfo.Validate(); err != nil {
		return err
	}

	return nil
}
