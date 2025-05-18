package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/LydiaTrack/ground/pkg/auth/types"
	"google.golang.org/api/idtoken"
)

type GoogleProvider struct {
	clientID     string
	clientSecret string
	redirectURI  string
	httpClient   *http.Client
}

func NewGoogleProvider(clientID, clientSecret, redirectURI string) *GoogleProvider {
	return &GoogleProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// GetUserInfo validates the token and returns user information
func (p *GoogleProvider) GetUserInfo(token string) (*types.OAuthUserInfo, error) {
	// Try direct access to userinfo endpoint first (works with access tokens)
	// This is often more reliable with tokens from the frontend
	userInfo, err := p.getUserInfoFromAPI(token)
	if err == nil {
		return userInfo, nil
	}

	// If that fails, try to validate as an ID token if it has the right format
	if len(strings.Split(token, ".")) == 3 {
		payload, err := idtoken.Validate(context.Background(), token, p.clientID)
		if err == nil {
			// Token is a valid ID token, extract user info from payload claims
			userInfo := &types.OAuthUserInfo{
				ProviderID: payload.Subject,
			}

			// Safely extract claims
			if email, ok := payload.Claims["email"].(string); ok {
				userInfo.Email = email
			}
			if name, ok := payload.Claims["name"].(string); ok {
				userInfo.Name = name
			}
			if givenName, ok := payload.Claims["given_name"].(string); ok {
				userInfo.FirstName = givenName
			}
			if familyName, ok := payload.Claims["family_name"].(string); ok {
				userInfo.LastName = familyName
			}
			if picture, ok := payload.Claims["picture"].(string); ok {
				userInfo.Picture = picture
			}

			return userInfo, nil
		}
	}

	// If all methods fail, return the original error
	return nil, fmt.Errorf("failed to get user info from token")
}

// getUserInfoFromAPI attempts to get user information using the Google userinfo API
func (p *GoogleProvider) getUserInfoFromAPI(token string) (*types.OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status code %d", resp.StatusCode)
	}

	var userInfo struct {
		ID         string `json:"id"`
		Email      string `json:"email"`
		Name       string `json:"name"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &types.OAuthUserInfo{
		ProviderID: userInfo.ID,
		Email:      userInfo.Email,
		Name:       userInfo.Name,
		FirstName:  userInfo.GivenName,
		LastName:   userInfo.FamilyName,
		Picture:    userInfo.Picture,
	}, nil
}
