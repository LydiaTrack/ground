package api

import (
	"github.com/LydiaTrack/lydia-base/internal/handlers"
	"github.com/LydiaTrack/lydia-base/internal/log"
	"github.com/LydiaTrack/lydia-base/pkg/middlewares"
	"github.com/LydiaTrack/lydia-base/pkg/service_initializer"
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
