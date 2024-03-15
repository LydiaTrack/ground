package user

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
	"net/mail"
	"time"
)

// Model user main model
type Model struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Username    string        `json:"username" bson:"username"`
	Password    string        `json:"-" bson:"password"`
	PersonInfo  `json:"personInfo" bson:"personInfo"`
	CreatedDate time.Time       `json:"createdDate" bson:"createdDate"`
	Version     int             `json:"version" bson:"version"`
	RoleIds     []bson.ObjectId `json:"roleIds" bson:"roleIds"`
}

func NewUser(id string, username string, password string, personInfo PersonInfo, createdDate time.Time, version int) Model {
	return Model{
		ID:          bson.ObjectIdHex(id),
		Username:    username,
		Password:    password,
		PersonInfo:  personInfo,
		CreatedDate: createdDate,
		Version:     version,
		RoleIds:     make([]bson.ObjectId, 0),
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

type PersonInfo struct {
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	Email       string             `json:"email,omitempty"`
	BirthDate   primitive.DateTime `json:"birthDate,omitempty"`
	Address     string             `json:"address,omitempty"`
	PhoneNumber `json:"phoneNumber,omitempty"`
}

func (p PersonInfo) Validate() error {
	if len(p.FirstName) == 0 {
		return errors.New("first name is required")
	}

	if len(p.LastName) == 0 {
		return errors.New("last name is required")
	}

	if len(p.Email) > 0 {
		if _, err := mail.ParseAddress(p.Email); err != nil {
			return err
		}
	}

	/*if err := p.PhoneNumber.Validate(); err != nil {
		return err
	}*/

	return nil
}

type PhoneNumber struct {
	AreaCode    string `json:"areaCode"`
	Number      string `json:"number"`
	CountryCode string `json:"countryCode"`
}

// Validate validates a phone number
func (p PhoneNumber) Validate() error {
	//TODO: More detailed validation can be done here
	if len(p.AreaCode) == 0 {
		return errors.New("area code is required")
	}

	if len(p.Number) == 0 {
		return errors.New("number is required")
	}

	if len(p.CountryCode) == 0 {
		return errors.New("country code is required")
	}

	return nil
}
