package middlewares

import (
	"github.com/LydiaTrack/lydia-base/internal/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware is a middleware for JWT authentication
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.ExtractTokenFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		err = jwt.IsTokenValid(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
