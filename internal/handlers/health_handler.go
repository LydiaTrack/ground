package handlers

import (
	"net/http"

	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

// GetHealth godoc
// @Summary Get health
// @Description get health.
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h HealthHandler) GetHealth(c *gin.Context) {
	health := service.GetApplicationHealth()
	c.JSON(http.StatusOK, health)
}
