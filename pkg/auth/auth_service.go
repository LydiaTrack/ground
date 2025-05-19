package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/LydiaTrack/ground/pkg/log"

	"github.com/LydiaTrack/ground/pkg/auth/providers"
	"github.com/LydiaTrack/ground/pkg/auth/types"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/session"
	"github.com/LydiaTrack/ground/pkg/domain/user"

	"github.com/LydiaTrack/ground/internal/jwt"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Create(command user.CreateUserCommand, authContext PermissionContext) (user.Model, error)
	ExistsByUsername(username string, authContext PermissionContext) (bool, error)
	ExistsByEmail(email string, authContext PermissionContext) (bool, error)
	VerifyUser(username, password string, authContext PermissionContext) (user.Model, error)
	Get(id string, authContext PermissionContext) (user.Model, error)
	GetByEmail(email string, authContext PermissionContext) (user.Model, error)
	Update(id string, command user.UpdateUserCommand, authContext PermissionContext) (user.Model, error)
}

type SessionService interface {
	DeleteSessionByUser(userID string) error
	CreateSession(command session.CreateSessionCommand) (session.InfoModel, error)
	GetUserSession(userID string) (session.InfoModel, error)
	GetSessionByRefreshToken(refreshToken string) (session.InfoModel, error)
}

type Service struct {
	userService    UserService
	sessionService SessionService
	oauthProviders map[string]types.OAuthProvider
}

type Response struct {
	jwt.TokenPair
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthService(userService UserService, sessionService SessionService) *Service {
	oauthProviders := make(map[string]types.OAuthProvider)

	// Initialize Google provider if credentials are available
	if googleClientID := os.Getenv("GOOGLE_CLIENT_ID"); googleClientID != "" {
		googleProvider := providers.NewGoogleProvider(
			googleClientID,
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_REDIRECT_URI"),
		)
		oauthProviders[types.GoogleProvider] = googleProvider
	}

	// Initialize Apple provider if credentials are available
	if appleClientID := os.Getenv("APPLE_CLIENT_ID"); appleClientID != "" {
		appleProvider := providers.NewAppleProvider(
			appleClientID,
			os.Getenv("APPLE_TEAM_ID"),
			os.Getenv("APPLE_KEY_ID"),
			os.Getenv("APPLE_PRIVATE_KEY"),
			os.Getenv("APPLE_REDIRECT_URI"),
		)
		oauthProviders[types.AppleProvider] = appleProvider
	}

	return &Service{
		userService:    userService,
		sessionService: sessionService,
		oauthProviders: oauthProviders,
	}
}

// Login is a function that handles the login process
func (s Service) Login(request Request) (Response, error) {
	// Check if user exists
	exists, err := s.userService.ExistsByUsername(request.Username, CreateAdminAuthContext())
	if err != nil {
		log.Log("Error checking if user exists", err)
		return Response{}, constants.ErrorInternalServerError
	}
	if !exists {
		log.Log("User does not exist", request.Username)
		return Response{}, constants.ErrorNotFound
	}

	// Check if password is correct
	userModel, err := s.userService.VerifyUser(request.Username, request.Password, PermissionContext{
		Permissions: []Permission{AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		log.Log("Error verifying user", err)
		return Response{}, err
	}

	// Generate token
	tokenPair, err := jwt.GenerateTokenPair(userModel.ID)
	if err != nil {
		log.Log("Error generating token pair", err)
		return Response{}, constants.ErrorInternalServerError
	}
	log.Log("Token pair", tokenPair)

	err = s.SetSession(userModel.ID.Hex(), tokenPair)
	if err != nil {
		return Response{}, constants.ErrorInternalServerError
	}

	return Response{
		tokenPair,
	}, nil
}

// SignUp is a function that handles the signup process, creates a new user from the given request
func (s Service) SignUp(cmd user.CreateUserCommand) (user.Model, error) {
	// Check if user exists
	exists, err := s.userService.ExistsByUsername(cmd.Username, CreateAdminAuthContext())
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}
	if exists {
		return user.Model{}, constants.ErrorConflict
	}

	// Create user
	userResponse, err := s.userService.Create(cmd, PermissionContext{Permissions: []Permission{AdminPermission}, UserID: nil})
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return userResponse, nil
}

// SetSession is a function that sets the session with the given user id and token pair
func (s Service) SetSession(userID string, tokenPair jwt.TokenPair) error {
	// Start a session
	refreshTokenLifespan, err := strconv.Atoi(os.Getenv(jwt.RefreshExpirationKey))
	if err != nil {
		return constants.ErrorInternalServerError
	}

	// If there is a session for the user, delete it
	err = s.sessionService.DeleteSessionByUser(userID)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	// Save refresh token with expire time
	createSessionCmd := session.CreateSessionCommand{
		UserID:       userID,
		ExpireTime:   time.Now().Add(time.Hour * time.Duration(refreshTokenLifespan)).Unix(),
		RefreshToken: tokenPair.RefreshToken,
	}
	_, err = s.sessionService.CreateSession(createSessionCmd)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	return nil
}

// GetCurrentUser is a function that returns the current user
func (s Service) GetCurrentUser(c *gin.Context) (user.Model, error) {
	userID, err := jwt.ExtractUserIDFromContext(c)
	if err != nil {
		return user.Model{}, constants.ErrorUnauthorized
	}

	// TODO: Maybe we should (or must) use GetSelfUser instead of Get, but I'm not sure.
	userModel, err := s.userService.Get(userID, CreateAdminAuthContext())
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return userModel, nil
}

// RefreshTokenPair is a function that refreshes the token pair
func (s Service) RefreshTokenPair(c *gin.Context) (jwt.TokenPair, error) {
	// Get the refresh token from the request body
	var refreshTokenRequest RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshTokenRequest); err != nil {
		return jwt.TokenPair{}, constants.ErrorInternalServerError
	}

	if refreshTokenRequest.RefreshToken == "" {
		return jwt.TokenPair{}, constants.ErrorUnauthorized
	}

	// Get the session by user id
	sessionInfo, err := s.sessionService.GetSessionByRefreshToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		return jwt.TokenPair{}, constants.ErrorInternalServerError
	}

	// Check if the refresh token is valid
	if sessionInfo.RefreshToken != refreshTokenRequest.RefreshToken {
		return jwt.TokenPair{}, constants.ErrorUnauthorized
	}

	// Now that we know the token is valid, we can extract the user id from it
	tokenPair, err := jwt.GenerateTokenPair(sessionInfo.UserID)
	if err != nil {
		return jwt.TokenPair{}, constants.ErrorInternalServerError
	}

	err = s.SetSession(sessionInfo.UserID.Hex(), tokenPair)
	if err != nil {
		return jwt.TokenPair{}, constants.ErrorInternalServerError
	}

	return tokenPair, nil
}

// HasPermission Checks if Permissions contains Permission
// It checks for the following cases:
// 1. */*
// 2. */Action
// 3. Domain/*
// 4. Domain/Action
func HasPermission(Permissions []Permission, Permission Permission) bool {
	// Check if there is a */*
	for _, permission := range Permissions {
		if permission.Domain == "*" && permission.Action == "*" {
			return true
		}
	}

	// Check if there is a */Action
	for _, permission := range Permissions {
		if permission.Domain == "*" && permission.Action == Permission.Action {
			return true
		}
	}

	// Check if there is a Domain/*
	for _, permission := range Permissions {
		if permission.Domain == Permission.Domain && permission.Action == "*" {
			return true
		}
	}

	// Check if there is a Domain/Action
	for _, permission := range Permissions {
		if permission.Domain == Permission.Domain && permission.Action == Permission.Action {
			return true
		}
	}

	return false
}

// CheckPermission checks if Permissions contains Permission
func CheckPermission(Permissions []Permission, Permission Permission) error {
	if !HasPermission(Permissions, Permission) {
		return constants.ErrorPermissionDenied
	}

	return nil
}

// OAuthLogin handles OAuth authentication
func (s Service) OAuthLogin(provider string, token string) (Response, error) {
	// Get the OAuth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		return Response{}, constants.ErrorBadRequest
	}

	// Get user info from the provider using the token
	userInfo, err := oauthProvider.GetUserInfo(token)
	if err != nil {
		return Response{}, constants.ErrorUnauthorized
	}

	// Check if user exists with this OAuth provider
	exists, err := s.userService.ExistsByEmail(userInfo.Email, CreateAdminAuthContext())
	if err != nil {
		return Response{}, constants.ErrorInternalServerError
	}

	oauthInfo := user.OAuthInfo{
		ProviderID:     userInfo.ProviderID,
		Email:          userInfo.Email,
		AccessToken:    token,
		RefreshToken:   "", // We no longer track refresh tokens separately
		TokenExpiry:    time.Now().Add(types.DefaultTokenExpiry),
		LastActiveDate: time.Now(),
	}
	var userModel user.Model
	if !exists {
		// Create new user
		createCmd := user.CreateUserCommand{
			Username: userInfo.Email, // Use email as username for OAuth users
			ContactInfo: user.ContactInfo{
				Email: userInfo.Email,
			},
			PersonInfo: &user.PersonInfo{
				FirstName: userInfo.FirstName,
				LastName:  userInfo.LastName,
			},
			// Picture is a link.
			Avatar:    userInfo.Picture,
			OAuthInfo: &oauthInfo,
		}
		userModel, err = s.userService.Create(createCmd, CreateAdminAuthContext())
		if err != nil {
			return Response{}, err
		}
	} else {
		// Get existing user
		userModel, err = s.userService.GetByEmail(userInfo.Email, CreateAdminAuthContext())
		if err != nil {
			return Response{}, err
		}
		// Update OAuth provider info
		if userModel.OAuthInfo == nil {
			userModel.OAuthInfo = &oauthInfo
		}
	}

	// Update user with OAuth info
	anyPropChanged := userModel.Avatar != userInfo.Picture ||
		userModel.PersonInfo.FirstName != userInfo.FirstName ||
		userModel.PersonInfo.LastName != userInfo.LastName ||
		userModel.ContactInfo.Email != userInfo.Email
	if anyPropChanged {
		updateCmd := user.UpdateUserCommand{
			Avatar: userInfo.Picture,
			PersonInfo: &user.PersonInfo{
				FirstName: userInfo.FirstName,
				LastName:  userInfo.LastName,
			},
			ContactInfo: &user.ContactInfo{
				Email: userInfo.Email,
			},
		}
		_, err = s.userService.Update(userModel.ID.Hex(), updateCmd, CreateAdminAuthContext())
		if err != nil {
			return Response{}, err
		}
	}

	// Generate JWT token
	tokenPair, err := jwt.GenerateTokenPair(userModel.ID)
	if err != nil {
		return Response{}, err
	}

	// Set session
	err = s.SetSession(userModel.ID.Hex(), tokenPair)
	if err != nil {
		return Response{}, err
	}

	return Response{tokenPair}, nil
}

// IsOAuthProviderEnabled checks if a specific OAuth provider is enabled
func (s Service) IsOAuthProviderEnabled(provider string) bool {
	_, exists := s.oauthProviders[provider]
	return exists
}
