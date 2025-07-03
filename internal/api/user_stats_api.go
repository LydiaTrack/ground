package api

import (
	"github.com/LydiaTrack/ground/internal/handlers"
	"github.com/LydiaTrack/ground/internal/repository"
	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/middlewares"
	"github.com/LydiaTrack/ground/pkg/service_initializer"
	"github.com/gin-gonic/gin"
)

func InitUserStats(r *gin.Engine) {
	userStatsService := service.NewUserStatsService(repository.GetUserStatsMongoRepository())
	userStatsHandler := handlers.NewUserStatsHandler(userStatsService, *service_initializer.GetServices().AuthService)

	routerGroup := r.Group("/user-stats")
	routerGroup.Use(middlewares.JwtAuthMiddleware()).
		GET("/me", userStatsHandler.GetSelfStats).
		GET("/me/task-streak", userStatsHandler.GetSelfTaskStreak)
}
