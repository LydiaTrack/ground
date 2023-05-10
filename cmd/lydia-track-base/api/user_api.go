package api

import (
	"github.com/gin-gonic/gin"
	"lydia-track-base/cmd/lydia-track-base/handlers"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
)

// InitUser initializes user routes
func InitUser(r *gin.Engine) {
	userRepository := repository.NewUserMongoRepository()
	userService := service.NewUserService(userRepository)

	userHandler := handlers.NewUserHandler(userService)

	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/users", userHandler.CreateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
}
