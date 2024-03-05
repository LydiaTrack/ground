package api

import (
	"github.com/Lydia/lydia-base/handlers"
	"github.com/gin-gonic/gin"
)

// InitHealth initializes health routes
func InitHealth(r *gin.Engine) {
	healthHandler := handlers.NewHealthHandler()

	r.GET("/health", healthHandler.GetHealth)
}
