package api

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/handlers"
)

// InitHealth initializes health routes
func InitHealth(r *gin.Engine) {
	healthHandler := handlers.NewHealthHandler()

	r.GET("/health", healthHandler.GetHealth)
}
