package service_initializer

import (
	"github.com/LydiaTrack/ground/internal/repository"
	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/auth"
)

type Services struct {
	AuthService          *auth.Service
	RoleService          *service.RoleService
	SessionService       *service.SessionService
	UserService          *service.UserService
	ResetPasswordService *service.ResetPasswordService
	FeedbackService      *service.FeedbackService
}

var services Services

// InitializeServices initializes all services and assigns them to the gin Engine
func InitializeServices() {
	services.RoleService = service.NewRoleService(repository.GetRoleMongoRepository())
	services.UserService = service.NewUserService(repository.GetUserMongoRepository(), *services.RoleService)
	services.SessionService = service.NewSessionService(repository.GetSessionRepository(), *services.UserService)
	services.AuthService = auth.NewAuthService(*services.UserService, *services.SessionService)
	services.ResetPasswordService = service.NewResetPasswordService(repository.GetResetPasswordRepository(), *services.UserService)
	services.FeedbackService = service.NewFeedbackService(repository.GetFeedbackRepository(), *services.UserService)
}

// GetServices returns the services.
func GetServices() Services {
	return services
}
