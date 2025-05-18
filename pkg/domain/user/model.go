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
	RoleIDs                  *[]primitive.ObjectID  `json:"roleIDs" bson:"roleIds"`
	Properties               map[string]interface{} `json:"properties" bson:"properties"`
	OAuthInfo                *OAuthInfo             `json:"OAuthInfo,omitempty" bson:"OAuthInfo,omitempty"`
}

type Option func(*Model) error

func NewUser(opts ...Option) (*Model, error) {
	u := &Model{
		ID:          primitive.NewObjectID(),
		CreatedDate: time.Now(),
		Version:     1,
		RoleIDs:     &[]primitive.ObjectID{},
	}

	for _, opt := range opts {
		if err := opt(u); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func WithUsername(username string) Option {
	return func(u *Model) error {
		u.Username = username
		return nil
	}
}

func WithPassword(password string) Option {
	return func(u *Model) error {
		u.Password = password
		return nil
	}
}

func WithAvatar(avatar string) Option {
	return func(u *Model) error {
		u.Avatar = avatar
		return nil
	}
}

func WithPersonInfo(personInfo *PersonInfo) Option {
	return func(u *Model) error {
		if personInfo == nil {
			// Skip setting the PersonInfo if it's nil
			return nil
		}
		u.PersonInfo = personInfo
		return nil
	}
}

func WithContactInfo(contactInfo ContactInfo) Option {
	return func(u *Model) error {
		u.ContactInfo = contactInfo
		return nil
	}
}

func WithProperties(properties map[string]interface{}) Option {
	return func(u *Model) error {
		u.Properties = properties
		return nil
	}
}

func WithRoleIDs(roleIDs *[]primitive.ObjectID) Option {
	return func(u *Model) error {
		u.RoleIDs = roleIDs
		return nil
	}
}

func WithOAuthInfo(oauthInfo *OAuthInfo) Option {
	return func(u *Model) error {
		u.OAuthInfo = oauthInfo
		return nil
	}
}

func (u Model) Validate() error {
	if u.Password == "" && u.OAuthInfo == nil {
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
		if err := utils.ValidateUserAvatar(u.Avatar); err != nil {
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

// OAuthInfo represents OAuth provider information for a user
type OAuthInfo struct {
	ProviderID    string    `json:"providerId" bson:"providerId"`
	Email         string    `json:"email" bson:"email"`
	AccessToken   string    `json:"-" bson:"accessToken"`
	RefreshToken  string    `json:"-" bson:"refreshToken"`
	TokenExpiry   time.Time `json:"tokenExpiry" bson:"tokenExpiry"`
	LastLoginDate time.Time `json:"lastLoginDate" bson:"lastLoginDate"`
}
