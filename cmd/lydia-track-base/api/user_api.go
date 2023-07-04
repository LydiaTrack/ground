package api

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/handlers"
	"lydia-track-base/internal/middlewares"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
)

// InitUser initializes user routes
func InitUser(r *gin.Engine) {

	userRepository := repository.GetUserRepository()
	userService := service.NewUserService(userRepository)

	userHandler := handlers.NewUserHandler(userService)

	authorizedPath := r.Group("/users")
	authorizedPath.Use(middlewares.JwtAuthMiddleware()).
		POST("", userHandler.CreateUser).
		GET("/:id", userHandler.GetUser).
		DELETE("/:id", userHandler.DeleteUser)
}
