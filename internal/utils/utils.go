package utils

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/gin-gonic/gin"
)

func CreateAuthContext(c *gin.Context, authService auth.Service, userService service.UserService) (auth.AuthContext, error) {
	currentUser, err := authService.GetCurrentUser(c)
	if err != nil {
		return auth.AuthContext{}, err
	}
	currentUserPermissions, err := userService.GetUserPermissionList(currentUser.ID)
	if err != nil {
		return auth.AuthContext{}, err
	}

	return auth.AuthContext{
		Permissions: currentUserPermissions,
		UserId:      currentUser.ID,
	}, nil
}
