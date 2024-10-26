package api

import (
	"github.com/LydiaTrack/ground/internal/handlers"
	"github.com/LydiaTrack/ground/pkg/service_initializer"
	"github.com/gin-gonic/gin"
)

func InitResetPassword(r *gin.Engine, services service_initializer.Services) {
	resetPasswordHandler := handlers.NewResetPasswordHandler(*services.UserService, *services.ResetPasswordService)

	routeGroup := r.Group("/reset-password")
	routeGroup.POST("", resetPasswordHandler.ResetPassword)
	routeGroup.POST("/verify-code", resetPasswordHandler.VerifyResetPasswordCode)
	routeGroup.POST("/send-email", resetPasswordHandler.SendResetPasswordCode)
}
