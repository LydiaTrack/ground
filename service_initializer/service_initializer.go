package service_initializer

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/internal/repository"
	"github.com/LydiaTrack/lydia-base/internal/service"
)

type Services struct {
	AuthService    *auth.Service
	RoleService    *service.RoleService
	SessionService *service.SessionService
	UserService    *service.UserService
}

var services Services

// InitializeServices initializes all services and assigns them to the gin Engine
func InitializeServices() {
	services.UserService = service.NewUserService(repository.GetUserRepository())
	services.RoleService = service.NewRoleService(repository.GetRoleRepository())
	services.SessionService = service.NewSessionService(repository.GetSessionRepository(), *services.UserService)
	services.AuthService = auth.NewAuthService(*services.UserService, *services.SessionService)
}

// GetServices returns the services.
func GetServices() Services {
	return services
}
