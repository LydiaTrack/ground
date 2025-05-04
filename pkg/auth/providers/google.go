package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LydiaTrack/ground/pkg/auth/types"
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

func (p *GoogleProvider) ValidateToken(token string) (*types.OAuthToken, error) {
	// Verify token with Google's tokeninfo endpoint
	resp, err := p.httpClient.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	var tokenInfo struct {
		Email        string `json:"email"`
		ExpiresIn    int64  `json:"expires_in"`
		AccessType   string `json:"access_type"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, err
	}

	return &types.OAuthToken{
		AccessToken:  token,
		RefreshToken: tokenInfo.RefreshToken,
		ExpiresIn:    tokenInfo.ExpiresIn,
		TokenType:    tokenInfo.TokenType,
		Expiry:       time.Now().Add(time.Duration(tokenInfo.ExpiresIn) * time.Second),
	}, nil
}

func (p *GoogleProvider) GetUserInfo(token *types.OAuthToken) (*types.OAuthUserInfo, error) {
	// Get user info from Google's userinfo endpoint
	req, err := http.NewRequestWithContext(context.Background(), "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info")
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
