package api

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/handlers"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
)

// InitAuth initializes auth routes
func InitAuth(r *gin.Engine) {
	userService := service.NewUserService(repository.GetUserRepository())
	sessionService := service.NewSessionService(repository.GetSessionRepository(), userService)
	authService := service.NewAuthService(userService, sessionService)

	authHandler := handlers.NewAuthHandler(authService)

	routeGroup := r.Group("/auth")
	routeGroup.POST("/login", authHandler.Login)
	routeGroup.GET("/currentUser", authHandler.GetCurrentUser)
	routeGroup.POST("/refreshToken", authHandler.RefreshToken)
}
