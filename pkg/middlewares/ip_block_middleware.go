package middlewares

import (
	"net/http"

	"github.com/LydiaTrack/lydia-base/internal/blocker"
	"github.com/gin-gonic/gin"
)

// IPBlockMiddleware checks if the IP is blocked
func IPBlockMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		method := c.Request.Method
		endpoint := c.FullPath()

		if blocker.GlobalBlocker.IsBlocked(ip, method, endpoint) {
			c.JSON(http.StatusForbidden, gin.H{"message": "You are temporarily blocked"})
			c.Abort()
			return
		}

		c.Next()
	}
}
