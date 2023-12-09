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
	userService := service.NewUserService(repository.GetUserRepository())
	sessionService := service.NewSessionService(repository.GetSessionRepository(), userService)
	authService := service.NewAuthService(userService, sessionService)

	roleHandler := handlers.NewRoleHandler(roleService, authService, userService)

	routerGroup := r.Group("/roles")
	routerGroup.GET("/:id", roleHandler.GetRole)
	routerGroup.POST("", roleHandler.CreateRole)
	routerGroup.DELETE("/:id", roleHandler.DeleteRole)

	utils.Log("Role routes initialized")
}
