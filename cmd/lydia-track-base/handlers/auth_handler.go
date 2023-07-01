package handlers

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/internal/auth"
)

type AuthHandler struct {
	authService auth.AuthService
}

func NewAuthHandler(authService auth.AuthService) AuthHandler {
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.Login(loginCommand)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description get current user.
// @Tags auth
// @Accept */*
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /current-user [get]
func (h AuthHandler) GetCurrentUser(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}
