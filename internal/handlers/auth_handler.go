package handlers

import (
	"github.com/LydiaTrack/ground/internal/service"
	"net/http"
	"time"

	"github.com/LydiaTrack/ground/internal/blocker"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/auth/types"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/log"
	"github.com/LydiaTrack/ground/pkg/service_initializer"
	"github.com/LydiaTrack/ground/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.Service
}

var userStatsService *service.UserStatsService

func NewAuthHandler(authService auth.Service) AuthHandler {
	userStatsService = service_initializer.GetServices().UserStatsService
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

	// Get the current user from the auth service to record login stats
	userStatsService := service_initializer.GetServices().UserStatsService
	if userStatsService != nil {
		// For login stats recording, we need the user's ID
		// We can get it by querying the user service with the login username
		userService := service_initializer.GetServices().UserService
		adminContext := auth.PermissionContext{
			Permissions: []auth.Permission{auth.AdminPermission},
			UserID:      nil,
		}

		userModel, err := userService.GetByUsername(loginCommand.Username, adminContext)
		if err == nil {
			// Record the login in user stats
			if recErr := userStatsService.RecordLogin(userModel.ID, adminContext); recErr != nil {
				// Log the error but don't fail the login process
				log.Log("Error recording login stats: %v", recErr)
			}
		}
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

	if userStatsService != nil {
		adminContext := auth.PermissionContext{
			Permissions: []auth.Permission{auth.AdminPermission},
			UserID:      nil,
		}

		// The stats should already be created as part of user creation, but let's verify
		_, err := userStatsService.GetUserStats(response.ID, adminContext)
		if err != nil {
			// If stats weren't created, create them now
			if createErr := userStatsService.CreateUserStats(response.ID, response.Username); createErr != nil {
				log.Log("Failed to create stats for new user: %v", createErr)
			}
		}
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
	response, err := h.authService.RefreshTokenPair(c)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	// Record the login in user stats
	if userStatsService != nil {
		adminContext := auth.PermissionContext{
			Permissions: []auth.Permission{auth.AdminPermission},
			UserID:      nil,
		}

		if recErr := userStatsService.RecordLogin(response.UserID, adminContext); recErr != nil {
			// Log the error but don't fail the login process
			log.Log("Error recording login stats: %v", recErr)
		}
	}
	c.JSON(http.StatusOK, response)
}

// OAuthLogin godoc
// @Summary OAuth login
// @Description login with OAuth provider (Google or Apple).
// @Tags auth
// @Accept json
// @Produce json
// @Param provider path string true "OAuth provider (google or apple)"
// @Param token body string true "OAuth token"
// @Success 200 {object} auth.Response
// @Router /auth/oauth/{provider} [post]
func (h AuthHandler) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")
	if provider != types.GoogleProvider && provider != types.AppleProvider {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}

	var token struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.OAuthLogin(provider, token.Token)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	// Record the login in user stats
	if userStatsService != nil {
		adminContext := auth.PermissionContext{
			Permissions: []auth.Permission{auth.AdminPermission},
			UserID:      nil,
		}

		if recErr := userStatsService.RecordLogin(response.UserID, adminContext); recErr != nil {
			// Log the error but don't fail the login process
			log.Log("Error recording login stats: %v", recErr)
		}
	}

	c.JSON(http.StatusOK, response)
}
