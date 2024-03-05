package api

import (
	"github.com/Lydia/lydia-base/handlers"
	"github.com/Lydia/lydia-base/internal/repository"
	"github.com/Lydia/lydia-base/internal/service"
	"github.com/gin-gonic/gin"
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
