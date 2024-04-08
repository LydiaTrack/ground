package api

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/handlers"
	"github.com/LydiaTrack/lydia-base/internal/repository"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/LydiaTrack/lydia-base/internal/utils"
	"github.com/gin-gonic/gin"
)

// InitUser initializes role routes
func InitRole(r *gin.Engine) {
	roleService := service.NewRoleService(repository.GetRoleRepository())
	userService := service.NewUserService(repository.GetUserRepository())
	sessionService := service.NewSessionService(repository.GetSessionRepository(), userService)
	authService := auth.NewAuthService(userService, sessionService)

	roleHandler := handlers.NewRoleHandler(roleService, authService, userService)

	routerGroup := r.Group("/roles")
	routerGroup.GET("/:id", roleHandler.GetRole)
	routerGroup.POST("", roleHandler.CreateRole)
	routerGroup.DELETE("/:id", roleHandler.DeleteRole)

	utils.Log("Role routes initialized")
}
