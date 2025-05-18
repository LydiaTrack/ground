package providers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/LydiaTrack/ground/pkg/auth/types"
)

type AppleProvider struct {
	clientID    string
	teamID      string
	keyID       string
	privateKey  string
	redirectURI string
	httpClient  *http.Client
}

func NewAppleProvider(clientID, teamID, keyID, privateKey, redirectURI string) *AppleProvider {
	return &AppleProvider{
		clientID:    clientID,
		teamID:      teamID,
		keyID:       keyID,
		privateKey:  privateKey,
		redirectURI: redirectURI,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetUserInfo validates the token and returns user information
func (p *AppleProvider) GetUserInfo(token string) (*types.OAuthUserInfo, error) {
	// Parse the JWT token without verification (we're just extracting claims)
	// In production, you should validate the signature with Apple's public key
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Get the claims part (second part) of the token
	claimsPart := parts[1]

	// Add padding if needed
	if len(claimsPart)%4 != 0 {
		claimsPart += strings.Repeat("=", 4-len(claimsPart)%4)
	}

	// Decode base64
	claimsBytes, err := base64.URLEncoding.DecodeString(claimsPart)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token claims: %w", err)
	}

	// Parse the claims
	var claims map[string]interface{}
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %w", err)
	}

	// Extract required fields
	var sub, email string
	if val, ok := claims["sub"].(string); ok {
		sub = val
	}
	if val, ok := claims["email"].(string); ok {
		email = val
	}

	// If we couldn't get a subject ID, the token is invalid
	if sub == "" {
		return nil, fmt.Errorf("invalid token: missing subject ID")
	}

	return &types.OAuthUserInfo{
		ProviderID: sub,
		Email:      email,
		Name:       "", // Apple doesn't provide name by default
		FirstName:  "", // These would be populated if the user shared their name
		LastName:   "",
		Picture:    "",
	}, nil
}
