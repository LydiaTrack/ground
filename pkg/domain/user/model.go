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

// StatsDocument represents a flexible statistics document for a user
// Uses map[string]interface{} to allow dynamic fields without changing the Ground library structure
type StatsDocument map[string]interface{}

// StatsCore contains the core required fields for user stats
type StatsCore struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	Username    string             `json:"username" bson:"username"`
	CreatedDate time.Time          `json:"createdDate" bson:"createdDate"`
	UpdatedDate time.Time          `json:"updatedDate" bson:"updatedDate"`
}

// NewStatsDocument creates a new flexible stats document for a user
func NewStatsDocument(userID primitive.ObjectID, username string) StatsDocument {
	now := time.Now()
	stats := StatsDocument{
		"_id":         primitive.NewObjectID(),
		"userId":      userID,
		"username":    username,
		"createdDate": now,
		"updatedDate": now,
		// Core activity stats
		"totalLogins":     1,
		"lastActiveDate":  now,
		"activeDaysCount": 1,
		"dayAge":          0, // Initially 0 days old
	}
	return stats
}

// GetCoreFields extracts core stats fields from the document
func (s StatsDocument) GetCoreFields() StatsCore {
	return StatsCore{
		ID:          s.GetObjectID("_id"),
		UserID:      s.GetObjectID("userId"),
		Username:    s.GetString("username"),
		CreatedDate: s.GetTime("createdDate"),
		UpdatedDate: s.GetTime("updatedDate"),
	}
}

// Helper methods for type-safe field access
func (s StatsDocument) GetObjectID(key string) primitive.ObjectID {
	if val, ok := s[key].(primitive.ObjectID); ok {
		return val
	}
	return primitive.NilObjectID
}

func (s StatsDocument) GetString(key string) string {
	if val, ok := s[key].(string); ok {
		return val
	}
	return ""
}

func (s StatsDocument) GetInt(key string) int {
	if val, ok := s[key].(int); ok {
		return val
	}
	if val, ok := s[key].(int32); ok {
		return int(val)
	}
	if val, ok := s[key].(int64); ok {
		return int(val)
	}
	return 0
}

func (s StatsDocument) GetInt64(key string) int64 {
	if val, ok := s[key].(int64); ok {
		return val
	}
	if val, ok := s[key].(int32); ok {
		return int64(val)
	}
	if val, ok := s[key].(int); ok {
		return int64(val)
	}
	return 0
}

func (s StatsDocument) GetTime(key string) time.Time {
	if val, ok := s[key].(time.Time); ok {
		return val
	}
	if val, ok := s[key].(primitive.DateTime); ok {
		return val.Time()
	}
	return time.Time{}
}

// SetField sets a field value in the document
func (s StatsDocument) SetField(key string, value interface{}) {
	s[key] = value
}

// IncrementField increments a numeric field by the given amount
func (s StatsDocument) IncrementField(key string, increment interface{}) {
	switch inc := increment.(type) {
	case int:
		s[key] = s.GetInt(key) + inc
	case int32:
		s[key] = int32(s.GetInt(key)) + inc
	case int64:
		s[key] = s.GetInt64(key) + inc
	default:
		s[key] = increment
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
func (s StatsDocument) CalculateStatFields() {
	now := time.Now()
	createdDate := s.GetTime("createdDate")
	lastActiveDate := s.GetTime("lastActiveDate")
	activeDaysCount := s.GetInt("activeDaysCount")

	// Check if this is a new day compared to the last active date before updating it
	if !lastActiveDate.IsZero() && !isSameDay(lastActiveDate, now) {
		activeDaysCount++
		s.SetField("activeDaysCount", activeDaysCount)
	}

	// Update last active date
	s.SetField("lastActiveDate", now)

	// Calculate day age (days since signup)
	dayAge := int(now.Sub(createdDate).Hours() / 24)
	s.SetField("dayAge", dayAge)

	// Update the updated date
	s.SetField("updatedDate", now)
}

// Helper function to check if two times are on the same day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
