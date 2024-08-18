package utils

import (
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/constants"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService interface {
	GetUserPermissionList(userId primitive.ObjectID) ([]auth.Permission, error)
}

func CreateAuthContext(c *gin.Context, authService auth.Service, userService userService) (auth.PermissionContext, error) {
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
