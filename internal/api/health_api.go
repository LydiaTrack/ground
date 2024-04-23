package api

import (
	"github.com/LydiaTrack/lydia-base/internal/handlers"
	"github.com/gin-gonic/gin"
)

// InitHealth initializes health routes
func InitHealth(r *gin.Engine) {
	healthHandler := handlers.NewHealthHandler()

	r.GET("/health", healthHandler.GetHealth)
}
