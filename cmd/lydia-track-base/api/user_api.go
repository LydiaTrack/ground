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

	userHandler := handlers.NewUserHandler(userService)

	routerGroup := r.Group("/users")
	routerGroup.Use(middlewares.JwtAuthMiddleware()).
		POST("", userHandler.CreateUser).
		GET("/:id", userHandler.GetUser).
		DELETE("/:id", userHandler.DeleteUser)

	utils.Log("User routes initialized")
}
