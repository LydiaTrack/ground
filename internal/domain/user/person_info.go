package user

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/mail"
)

type PersonInfo struct {
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	Email       string             `json:"email,omitempty"`
	BirthDate   primitive.DateTime `json:"birth_date,omitempty"`
	Address     string             `json:"address,omitempty"`
	PhoneNumber `bson:"phone_number,omitempty"`
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
