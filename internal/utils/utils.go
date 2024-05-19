package utils

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/gin-gonic/gin"
)

func CreateAuthContext(c *gin.Context, authService auth.Service, userService service.UserService) (auth.PermissionContext, error) {
	currentUser, err := authService.GetCurrentUser(c)
	if err != nil {
		return auth.PermissionContext{}, err
	}
	currentUserPermissions, err := userService.GetUserPermissionList(currentUser.ID)
	if err != nil {
		return auth.PermissionContext{}, err
	}

	return auth.PermissionContext{
		Permissions: currentUserPermissions,
		UserId:      &currentUser.ID,
	}, nil
}
