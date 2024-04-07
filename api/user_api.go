package api

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/handlers"
	"github.com/LydiaTrack/lydia-base/internal/middlewares"
	"github.com/LydiaTrack/lydia-base/internal/repository"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/LydiaTrack/lydia-base/internal/utils"
	"github.com/gin-gonic/gin"
)

// InitUser initializes user routes
func InitUser(r *gin.Engine) {

	userService := service.NewUserService(repository.GetUserRepository())
	sessionService := service.NewSessionService(repository.GetSessionRepository(), userService)
	authService := auth.NewAuthService(userService, sessionService)

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
