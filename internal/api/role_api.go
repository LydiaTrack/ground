package api

import (
	"github.com/LydiaTrack/ground/internal/handlers"
	"github.com/LydiaTrack/ground/pkg/log"
	"github.com/LydiaTrack/ground/pkg/middlewares"
	"github.com/LydiaTrack/ground/pkg/service_initializer"
	"github.com/gin-gonic/gin"
)

// InitUser initializes role routes
func InitRole(r *gin.Engine, services service_initializer.Services) {

	roleHandler := handlers.NewRoleHandler(*services.RoleService, *services.AuthService, *services.UserService)

	routerGroup := r.Group("/roles")
	routerGroup.Use(middlewares.JwtAuthMiddleware())
	routerGroup.GET("", roleHandler.GetRoles)
	routerGroup.GET("/:id", roleHandler.GetRole)
	routerGroup.POST("", roleHandler.CreateRole)
	routerGroup.DELETE("/:id", roleHandler.DeleteRole)

	log.Log("Role routes initialized")
}
