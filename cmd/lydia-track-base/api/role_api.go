package api

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/handlers"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
	"lydia-track-base/internal/utils"
)

// InitUser initializes role routes
func InitRole(r *gin.Engine) {
	roleService := service.NewRoleService(repository.GetRoleRepository())
	roleHandler := handlers.NewRoleHandler(roleService)

	routerGroup := r.Group("/roles")
	routerGroup.GET("/:id", roleHandler.GetRole)
	routerGroup.POST("", roleHandler.CreateRole)
	routerGroup.DELETE("/:id", roleHandler.DeleteRole)

	utils.Log("Role routes initialized")
}
