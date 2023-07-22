package api

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/handlers"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
)

// InitUser initializes role routes
func InitRole(r *gin.Engine) {
	roleRepository := repository.NewRoleMongoRepository()
	roleService := service.NewRoleService(roleRepository)
	roleHandler := handlers.NewRoleHandler(roleService)

	
	r.GET("/roles/:id", roleHandler.GetRole)
	r.POST("/roles", roleHandler.CreateRole)
	r.DELETE("/roles/:id", roleHandler.DeleteRole)
}
