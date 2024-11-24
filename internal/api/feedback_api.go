package api

import (
	"github.com/LydiaTrack/ground/internal/handlers"
	"github.com/LydiaTrack/ground/pkg/service_initializer"
	"github.com/gin-gonic/gin"
)

func InitFeedback(r *gin.Engine, services service_initializer.Services) {
	feedbackHandler := handlers.NewFeedbackHandler(*services.FeedbackService)

	routeGroup := r.Group("/feedback")
	routeGroup.POST("", feedbackHandler.CreateFeedback)
	routeGroup.GET("/user/:userID", feedbackHandler.GetFeedbackByUser)
	routeGroup.PUT("/:feedbackID/status", feedbackHandler.UpdateFeedbackStatus)
}
