package middlewares

import (
	"lydia-track-base/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware is a middleware for JWT authentication
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.IsTokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
