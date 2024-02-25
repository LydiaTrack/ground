package middlewares

import (
	"github.com/Lydia/lydia-base/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware is a middleware for JWT authentication
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.ExtractTokenFromContext(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		err = utils.IsTokenValid(token)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
