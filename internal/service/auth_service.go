package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	auth2 "lydia-track-base/internal/auth"
	"lydia-track-base/internal/domain/auth"
	"lydia-track-base/internal/domain/session/commands"
	"lydia-track-base/internal/domain/user"
	"lydia-track-base/internal/service"
	"lydia-track-base/internal/utils"
	"os"
	"strconv"
	"time"
)

type Service struct {
	userService    UserService
	sessionService SessionService
}

type Response struct {
	auth2.TokenPair
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthService(userService UserService, sessionService SessionService) Service {
	return Service{
		userService:    userService,
		sessionService: sessionService,
	}
}

// Login is a function that handles the login process
func (s Service) Login(request Request) (Response, error) {
	// Check if user exists
	exists, err := s.userService.ExistsByUsername(request.Username, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return Response{}, err
	}
	if !exists {
		return Response{}, errors.New("user does not exist")
	}

	// Check if password is correct
	userModel, err := s.userService.VerifyUser(request.Username, request.Password, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return Response{}, err
	}

	// Generate token
	tokenPair, err := auth2.GenerateTokenPair(userModel.ID)
	if err != nil {
		return Response{}, err
	}

	err = s.SetSession(userModel.ID.Hex(), tokenPair)
	if err != nil {
		return Response{}, err
	}

	return Response{
		tokenPair,
	}, nil
}

// SetSession is a function that sets the session with the given user id and token pair
func (s Service) SetSession(userId string, tokenPair auth2.TokenPair) error {
	// Start a session
	refreshTokenLifespan, err := strconv.Atoi(os.Getenv(auth2.RefreshExpirationKey))
	if err != nil {
		return err
	}

	// If there is a session for the user, delete it
	err = s.sessionService.DeleteSessionByUser(userId)
	if err != nil {
		return err
	}

	// Save refresh token with expire time
	createSessionCmd := commands.CreateSessionCommand{
		UserId:       userId,
		ExpireTime:   time.Now().Add(time.Hour * time.Duration(refreshTokenLifespan)).Unix(),
		RefreshToken: tokenPair.RefreshToken,
	}
	_, err = s.sessionService.CreateSession(createSessionCmd)
	if err != nil {
		return err
	}

	return nil
}

// GetCurrentUser is a function that returns the current user
func (s Service) GetCurrentUser(c *gin.Context) (user.Model, error) {
	userId, err := auth2.ExtractUserIdFromContext(c)
	if err != nil {
		return user.Model{}, err
	}

	userModel, err := s.userService.GetUser(userId, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return user.Model{}, err
	}

	return userModel, nil
}

// RefreshTokenPair is a function that refreshes the token pair
func (s Service) RefreshTokenPair(c *gin.Context) (auth2.TokenPair, error) {
	// Get the refresh token from the request body
	var refreshTokenRequest RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshTokenRequest); err != nil {
		return auth2.TokenPair{}, err
	}

	// Get current user id
	currentUser, err := s.GetCurrentUser(c)
	if err != nil {
		return auth2.TokenPair{}, err
	}

	// Get the session by user id
	sessionInfo, err := s.sessionService.GetUserSession(currentUser.ID.Hex())
	if err != nil {
		return auth2.TokenPair{}, err
	}

	// Check if the refresh token is valid
	if sessionInfo.RefreshToken != refreshTokenRequest.RefreshToken {
		return auth2.TokenPair{}, errors.New("refresh token is invalid")
	}

	// Now that we know the token is valid, we can extract the user id from it
	tokenPair, err := auth2.GenerateTokenPair(currentUser.ID)
	if err != nil {
		return auth2.TokenPair{}, err
	}

	err = s.SetSession(currentUser.ID.Hex(), tokenPair)
	if err != nil {
		return auth2.TokenPair{}, err
	}

	return tokenPair, nil
}

// CheckPermission is a function that checks if Permissions contains Permission
// It checks for the following cases:
// 1. */*
// 2. */Action
// 3. Domain/*
// 4. Domain/Action
func CheckPermission(Permissions []auth.Permission, Permission auth.Permission) bool {
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
