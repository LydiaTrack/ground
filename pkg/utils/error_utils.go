package utils

import (
	"errors"
	"github.com/LydiaTrack/lydia-base/pkg/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EvaluateError(err error, c *gin.Context) {
	switch {
	case errors.Is(err, constants.ErrorUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, constants.ErrorPermissionDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, constants.ErrorBadRequest):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, constants.ErrorNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, constants.ErrorInternalServerError):
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case errors.Is(err, constants.ErrorConflict):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
