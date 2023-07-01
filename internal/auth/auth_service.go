package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"lydia-track-base/internal/domain"
)

type AuthService struct {
	userService UserService
}

type UserService interface {
	ExistsByUsername(username string) bool
	VerifyUser(username string, password string) (domain.UserModel, error)
	GetUser(id string) (domain.UserModel, error)
}

type Response struct {
	Token string `json:"token"`
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthService(userService UserService) AuthService {
	return AuthService{
		userService: userService,
	}
}

// Login is a function that handles the login process
func (a AuthService) Login(request Request) (Response, error) {
	// Check if user exists
	exists := a.userService.ExistsByUsername(request.Username)
	if !exists {
		return Response{}, errors.New("user does not exist")
	}

	// Check if password is correct
	user, err := a.userService.VerifyUser(request.Username, request.Password)
	if err != nil {
		return Response{}, err
	}

	// Generate token
	token, err := GenerateToken(user.ID)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Token: token,
	}, nil
}

// GetCurrentUser is a function that returns the current user
func (a AuthService) GetCurrentUser(c *gin.Context) (domain.UserModel, error) {
	userId, err := ExtractTokenID(c)
	if err != nil {
		return domain.UserModel{}, err
	}

	user, err := a.userService.GetUser(userId)
	if err != nil {
		return domain.UserModel{}, err
	}

	return user, nil
}
