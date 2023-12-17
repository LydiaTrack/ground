package api

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/handlers"
	"lydia-track-base/internal/middlewares"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
	"lydia-track-base/internal/utils"
)

// InitUser initializes user routes
func InitUser(r *gin.Engine) {

	userService := service.NewUserService(repository.GetUserRepository())
	sessionService := service.NewSessionService(repository.GetSessionRepository(), userService)
	authService := service.NewAuthService(userService, sessionService)

	userHandler := handlers.NewUserHandler(userService, authService)

	routerGroup := r.Group("/users")
	routerGroup.Use(middlewares.JwtAuthMiddleware()).
		POST("", userHandler.CreateUser).
		GET("/:id", userHandler.GetUser).
		DELETE("/:id", userHandler.DeleteUser).
		GET("/roles/:id", userHandler.GetUserRoles).
		POST("/roles", userHandler.AddRoleToUser).
		DELETE("/roles", userHandler.RemoveRoleFromUser)

	utils.Log("User routes initialized")
}
