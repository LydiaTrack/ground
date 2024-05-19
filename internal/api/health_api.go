package api

import (
	"github.com/LydiaTrack/lydia-base/internal/handlers"
	"github.com/LydiaTrack/lydia-base/middlewares"
	"github.com/gin-gonic/gin"
)

// InitHealth initializes health routes
func InitHealth(r *gin.Engine) {
	healthHandler := handlers.NewHealthHandler()

	routerGroup := r.Group("/health")
	routerGroup.Use(middlewares.JwtAuthMiddleware())
	routerGroup.GET("", healthHandler.GetHealth)
}
