package auth

import (
	"github.com/LydiaTrack/lydia-base/internal/log"
	"os"
	"strconv"
	"time"

	"github.com/LydiaTrack/lydia-base/pkg/constants"
	"github.com/LydiaTrack/lydia-base/pkg/domain/session"
	"github.com/LydiaTrack/lydia-base/pkg/domain/user"

	"github.com/LydiaTrack/lydia-base/internal/jwt"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(command user.CreateUserCommand, authContext PermissionContext) (user.Model, error)
	ExistsByUsername(username string) bool
	ExistsByEmail(email string) bool
	VerifyUser(username, password string, authContext PermissionContext) (user.Model, error)
	GetUser(id string, authContext PermissionContext) (user.Model, error)
}

type SessionService interface {
	DeleteSessionByUser(userId string) error
	CreateSession(command session.CreateSessionCommand) (session.InfoModel, error)
	GetUserSession(userId string) (session.InfoModel, error)
	GetSessionByRefreshToken(refreshToken string) (session.InfoModel, error)
}

type Service struct {
	userService    UserService
	sessionService SessionService
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
	return &Service{
		userService:    userService,
		sessionService: sessionService,
	}
}

// Login is a function that handles the login process
func (s Service) Login(request Request) (Response, error) {
	// Check if user exists
	exists := s.userService.ExistsByUsername(request.Username)
	if !exists {
		log.Log("User does not exist", request.Username)
		return Response{}, constants.ErrorNotFound
	}

	// Check if password is correct
	userModel, err := s.userService.VerifyUser(request.Username, request.Password, PermissionContext{
		Permissions: []Permission{AdminPermission},
		UserId:      nil,
	})
	if err != nil {
		log.Log("Error verifying user", err)
		return Response{}, constants.ErrorInternalServerError
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
	exists := s.userService.ExistsByUsername(cmd.Username)
	if exists {
		return user.Model{}, constants.ErrorConflict
	}

	// Create user
	userResponse, err := s.userService.CreateUser(cmd, PermissionContext{Permissions: []Permission{AdminPermission}, UserId: nil})
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return userResponse, nil
}

// SetSession is a function that sets the session with the given user id and token pair
func (s Service) SetSession(userId string, tokenPair jwt.TokenPair) error {
	// Start a session
	refreshTokenLifespan, err := strconv.Atoi(os.Getenv(jwt.RefreshExpirationKey))
	if err != nil {
		return constants.ErrorInternalServerError
	}

	// If there is a session for the user, delete it
	err = s.sessionService.DeleteSessionByUser(userId)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	// Save refresh token with expire time
	createSessionCmd := session.CreateSessionCommand{
		UserId:       userId,
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
	userId, err := jwt.ExtractUserIdFromContext(c)
	if err != nil {
		return user.Model{}, constants.ErrorUnauthorized
	}

	// TODO: Maybe we should (or must) use GetSelfUser instead of GetUser, but I'm not sure.
	userModel, err := s.userService.GetUser(userId, PermissionContext{
		Permissions: []Permission{AdminPermission},
		UserId:      nil,
	})
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
	tokenPair, err := jwt.GenerateTokenPair(sessionInfo.UserId)
	if err != nil {
		return jwt.TokenPair{}, constants.ErrorInternalServerError
	}

	err = s.SetSession(sessionInfo.UserId.Hex(), tokenPair)
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
