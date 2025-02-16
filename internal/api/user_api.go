package api

import (
	"github.com/LydiaTrack/ground/internal/handlers"
	"github.com/LydiaTrack/ground/pkg/log"
	"github.com/LydiaTrack/ground/pkg/middlewares"
	"github.com/LydiaTrack/ground/pkg/service_initializer"
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
		PUT("/:id", userHandler.UpdateUser).
		DELETE("/:id", userHandler.DeleteUser).
		GET("/roles/:id", userHandler.GetUserRoles).
		POST("/roles", userHandler.AddRoleToUser).
		DELETE("/roles", userHandler.RemoveRoleFromUser)
	checkUsernameGroup := r.Group("/users/checkUsername")
	checkUsernameGroup.GET("/:username", userHandler.CheckUsername)

	checkEmailGroup := r.Group("/users/checkEmail")
	checkEmailGroup.GET("/:email", userHandler.CheckEmail)

	selfRouterGroup := r.Group("/users-self")
	selfRouterGroup.Use(middlewares.JwtAuthMiddleware()).
		PUT("", userHandler.UpdateUserSelf).
		PUT("/password", userHandler.UpdateUserSelfPassword)

	log.Log("User routes initialized")
}
