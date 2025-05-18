package types

import "time"

const (
	// OAuth providers
	GoogleProvider = "google"
	AppleProvider  = "apple"

	// OAuth token expiration
	DefaultTokenExpiry = 1 * time.Hour
)

// OAuthToken represents the token response from OAuth providers
type OAuthToken struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresIn    int64     `json:"expiresIn"`
	TokenType    string    `json:"tokenType"`
	Expiry       time.Time `json:"expiry"`
}

// OAuthUserInfo represents the user information from OAuth providers
type OAuthUserInfo struct {
	ProviderID string `json:"providerId"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Picture    string `json:"picture"`
}

// OAuthProvider defines the interface for OAuth providers
type OAuthProvider interface {
	GetUserInfo(token string) (*OAuthUserInfo, error)
}
