package api

import (
	"github.com/LydiaTrack/lydia-base/internal/handlers"
	"github.com/LydiaTrack/lydia-base/internal/log"
	"github.com/LydiaTrack/lydia-base/pkg/middlewares"
	"github.com/LydiaTrack/lydia-base/pkg/service_initializer"
	"github.com/gin-gonic/gin"
)

// InitUser initializes user routes
func InitUser(r *gin.Engine, services service_initializer.Services) {

	userHandler := handlers.NewUserHandler(*services.UserService, *services.AuthService)

	routerGroup := r.Group("/users")
	routerGroup.Use(middlewares.JwtAuthMiddleware()).
		POST("", userHandler.CreateUser).
		GET("", userHandler.GetUsers).
		GET("/:id", userHandler.GetUser).
		DELETE("/:id", userHandler.DeleteUser).
		GET("/roles/:id", userHandler.GetUserRoles).
		POST("/roles", userHandler.AddRoleToUser).
		DELETE("/roles", userHandler.RemoveRoleFromUser)
	checkUsernameGroup := r.Group("/users/checkUsername")
	checkUsernameGroup.GET("/:username", userHandler.CheckUsername)

	log.Log("User routes initialized")
}
