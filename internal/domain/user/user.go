package user

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain/role"
	"time"
)

type Model struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Username    string        `json:"username" bson:"username"`
	Password    string        `json:"-" bson:"password"`
	PersonInfo  `json:"personInfo" bson:"personInfo"`
	CreatedDate time.Time    `json:"createdDate" bson:"createdDate"`
	Version     int          `json:"version" bson:"version"`
	Roles       []role.Model `json:"roles,omitempty" bson:"roles,omitempty"`
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
