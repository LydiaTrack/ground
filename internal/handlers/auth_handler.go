package handlers

import (
	"net/http"
	"time"

	"github.com/LydiaTrack/lydia-base/internal/blocker"
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/domain/user"
	"github.com/LydiaTrack/lydia-base/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) AuthHandler {
	return AuthHandler{authService: authService}
}

// Login godoc
// @Summary Login
// @Description login.
// @Tags auth
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func (h AuthHandler) Login(c *gin.Context) {
	var loginCommand auth.Request
	if err := c.ShouldBindJSON(&loginCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.Login(loginCommand)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.JSON(http.StatusOK, response)
}

// SignUp godoc
// @Summary Sign up
// @Description sign up.
// @Tags auth
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /signUp [post]
func (h AuthHandler) SignUp(c *gin.Context) {
	ip := c.ClientIP()
	method := c.Request.Method
	endpoint := c.FullPath()

	var createUserCommand user.CreateUserCommand
	if err := c.ShouldBindJSON(&createUserCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.SignUp(createUserCommand)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	// After user signs up, we should block the user from signing up again for a certain period of time
	// This is to prevent spam signups
	blocker.GlobalBlocker.Add(ip, method, endpoint, 5*time.Second)
	c.JSON(http.StatusOK, response)
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description get current user.
// @Tags auth
// @Accept */*
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /currentUser [get]
func (h AuthHandler) GetCurrentUser(c *gin.Context) {
	userModel, err := h.authService.GetCurrentUser(c)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.JSON(http.StatusOK, userModel)
}

// RefreshToken godoc
// @Summary Refresh token
// @Description refresh token.
// @Tags auth
// @Accept */*
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /refreshToken [get]
func (h AuthHandler) RefreshToken(c *gin.Context) {
	tokenPair, err := h.authService.RefreshTokenPair(c)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}
	c.JSON(http.StatusOK, tokenPair)
}
