package handlers

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService service.Service
}

func NewAuthHandler(authService service.Service) AuthHandler {
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
	var loginCommand service.Request
	if err := c.ShouldBindJSON(&loginCommand); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.Login(loginCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tokenPair)
}
