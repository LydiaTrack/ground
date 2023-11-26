package user

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain/role"
	"time"
)

type Model struct {
	ID          bson.ObjectId `bson:"_id"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	PersonInfo  `json:"personInfo"`
	CreatedDate time.Time    `json:"createdDate,omitempty"`
	Version     int          `json:"version,omitempty"`
	Roles       []role.Model `json:"roles,omitempty"`
}

func NewUser(id string, username string, password string, personInfo PersonInfo, createdDate time.Time, version int) Model {
	return Model{
		ID:          bson.ObjectIdHex(id),
		Username:    username,
		Password:    password,
		PersonInfo:  personInfo,
		CreatedDate: createdDate,
		Version:     version,
	}

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
