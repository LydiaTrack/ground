package user

import (
	"errors"
	"time"

	"github.com/LydiaTrack/ground/internal/utils"

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

// StatsModel represents the statistics for a user
type StatsModel struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	Username    string             `json:"username" bson:"username"`
	CreatedDate time.Time          `json:"createdDate" bson:"createdDate"`
	UpdatedDate time.Time          `json:"updatedDate" bson:"updatedDate"`

	// Activity stats
	TotalLogins     int       `json:"totalLogins" bson:"totalLogins"`
	LastActiveDate  time.Time `json:"lastActiveDate,omitempty" bson:"lastActiveDate,omitempty"`
	ActiveDaysCount int       `json:"activeDaysCount" bson:"activeDaysCount"`
	DayAge          int       `json:"dayAge" bson:"dayAge"` // Days since signup

	// Task stats
	TasksCreated   int `json:"tasksCreated" bson:"tasksCreated"`
	TasksCompleted int `json:"tasksCompleted" bson:"tasksCompleted"`

	// Note stats
	NotesCreated int `json:"notesCreated" bson:"notesCreated"`

	// Time tracking stats
	TotalTimeTracked int64 `json:"totalTimeTracked" bson:"totalTimeTracked"` // In seconds
	TimeEntryCount   int   `json:"timeEntryCount" bson:"timeEntryCount"`

	// Project stats
	ProjectsCreated int `json:"projectsCreated" bson:"projectsCreated"`
}

// NewStats creates a new stats model for a user
func NewStats(userID primitive.ObjectID, username string) *StatsModel {
	now := time.Now()
	return &StatsModel{
		ID:               primitive.NewObjectID(),
		UserID:           userID,
		Username:         username,
		CreatedDate:      now,
		UpdatedDate:      now,
		TotalLogins:      1,
		LastActiveDate:   now,
		ActiveDaysCount:  1,
		DayAge:           0, // Initially 0 days old
		TasksCreated:     0,
		TasksCompleted:   0,
		NotesCreated:     0,
		TotalTimeTracked: 0,
		TimeEntryCount:   0,
		ProjectsCreated:  0,
	}
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
	ProviderID     string    `json:"providerId" bson:"providerId"`
	Email          string    `json:"email" bson:"email"`
	AccessToken    string    `json:"-" bson:"accessToken"`
	RefreshToken   string    `json:"-" bson:"refreshToken"`
	TokenExpiry    time.Time `json:"tokenExpiry" bson:"tokenExpiry"`
	LastActiveDate time.Time `json:"lastActiveDate" bson:"lastActiveDate"`
}

// CalculateStatFields updates general stat fields that should be updated on every stat change
func (s *StatsModel) CalculateStatFields() {
	now := time.Now()

	// Update last active date
	s.LastActiveDate = now

	if s.LastActiveDate.Year() != now.Year() ||
		s.LastActiveDate.YearDay() != now.YearDay() {
		s.ActiveDaysCount++
	}

	// Calculate day age (days since signup)
	s.DayAge = int(now.Sub(s.CreatedDate).Hours() / 24)

	// Update the updated date
	s.UpdatedDate = now
}
