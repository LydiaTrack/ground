package user

import (
	"errors"
	"github.com/LydiaTrack/ground/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model user main model
type Model struct {
	ID                       primitive.ObjectID     `json:"id" bson:"_id"`
	Username                 string                 `json:"username" bson:"username"`
	Password                 string                 `json:"-" bson:"password"`
	Avatar                   string                 `json:"avatar,omitempty" bson:"avatar,omitempty"`
	PersonInfo               *PersonInfo            `json:"personInfo" bson:"personInfo"`
	ContactInfo              ContactInfo            `json:"contactInfo" bson:"contactInfo"`
	CreatedDate              time.Time              `json:"createdDate" bson:"createdDate"`
	Version                  int                    `json:"version" bson:"version"`
	LastSeenChangelogVersion string                 `json:"lastSeenChangelogVersion" bson:"lastSeenChangelogVersion"`
	RoleIds                  *[]primitive.ObjectID  `json:"roleIds" bson:"roleIds"`
	Properties               map[string]interface{} `json:"properties" bson:"properties"`
}

func NewUser(id string, username string, password string,
	personInfo *PersonInfo, contactInfo ContactInfo,
	createdDate time.Time, properties map[string]interface{}, version int) (Model, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Model{}, err
	}
	return Model{
		ID:          objID,
		Username:    username,
		Password:    password,
		PersonInfo:  personInfo,
		ContactInfo: contactInfo,
		CreatedDate: createdDate,
		Version:     version,
		RoleIds:     &[]primitive.ObjectID{},
		Properties:  properties,
		Avatar:      "",
	}, nil
}

func (u Model) Validate() error {
	if u.Password == "" {
		return errors.New("password is required")
	}

	if u.Username == "" {
		return errors.New("username is required")
	}

	if u.PersonInfo != nil {
		if err := u.PersonInfo.Validate(); err != nil {
			return err
		}
	}

	if u.Avatar != "" {
		if err := utils.ValidateBase64Image(u.Avatar); err != nil {
			return err
		}
	}

	return nil
}

type ContactInfo struct {
	Email       string       `json:"email,omitempty"`
	PhoneNumber *PhoneNumber `json:"phoneNumber,omitempty"`
}

type PersonInfo struct {
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	BirthDate primitive.DateTime `json:"birthDate,omitempty"`
}

func (p PersonInfo) Validate() error {
	if p.FirstName == "" {
		return errors.New("first name is required")
	}

	if p.LastName == "" {
		return errors.New("last name is required")
	}

	return nil
}

type PhoneNumber struct {
	AreaCode    string `json:"areaCode"`
	Number      string `json:"number"`
	CountryCode string `json:"countryCode"`
}

// Validate validates a phone number
func (p PhoneNumber) Validate() error {
	if p.AreaCode == "" {
		return errors.New("area code is required")
	}

	if p.Number == "" {
		return errors.New("number is required")
	}

	if p.CountryCode == "" {
		return errors.New("country code is required")
	}

	return nil
}
