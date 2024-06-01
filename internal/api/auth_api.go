package api

import (
	"github.com/LydiaTrack/lydia-base/internal/handlers"
	"github.com/LydiaTrack/lydia-base/pkg/service_initializer"
	"github.com/gin-gonic/gin"
)

// InitAuth initializes auth routes
func InitAuth(r *gin.Engine, services service_initializer.Services) {

	authHandler := handlers.NewAuthHandler(*services.AuthService)

	routeGroup := r.Group("/auth")
	routeGroup.POST("/login", authHandler.Login)
	routeGroup.GET("/currentUser", authHandler.GetCurrentUser)
	routeGroup.POST("/refreshToken", authHandler.RefreshToken)
}
