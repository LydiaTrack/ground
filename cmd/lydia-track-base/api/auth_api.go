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
	userRepository := repository.GetRepository()
	userService := service.NewUserService(userRepository)
	authService := auth.NewAuthService(userService)

	authHandler := handlers.NewAuthHandler(authService)

	r.POST("/login", authHandler.Login)
	r.GET("/currentUser", authHandler.GetCurrentUser)
}
