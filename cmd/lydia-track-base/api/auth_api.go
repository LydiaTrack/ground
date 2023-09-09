package api

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/handlers"
	"lydia-track-base/internal/auth"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
)

// InitAuth initializes auth routes
func InitAuth(r *gin.Engine) {
	userRepository := repository.GetUserRepository()
	userService := service.NewUserService(userRepository)
	sessionRepository := repository.GetSessionRepository()
	sessionService := service.NewSessionService(sessionRepository, userService)
	authService := auth.NewAuthService(userService, sessionService)

	authHandler := handlers.NewAuthHandler(authService)

	r.POST("/login", authHandler.Login)
	r.GET("/currentUser", authHandler.GetCurrentUser)
	r.POST("/refreshToken", authHandler.RefreshToken)
}
