package utils

import (
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/constants"
	"github.com/gin-gonic/gin"
)

func CreateAuthContext(c *gin.Context, authService auth.Service, userService service.UserService) (auth.PermissionContext, error) {
	currentUser, err := authService.GetCurrentUser(c)
	if err != nil {
		return auth.PermissionContext{}, constants.ErrorNotFound
	}
	currentUserPermissions, err := userService.GetUserPermissionList(currentUser.ID)
	if err != nil {
		return auth.PermissionContext{}, constants.ErrorNotFound
	}

	return auth.PermissionContext{
		Permissions: currentUserPermissions,
		UserId:      &currentUser.ID,
	}, nil
}
