package utils

import (
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService interface {
	GetPermissionList(userID primitive.ObjectID) ([]auth.Permission, error)
}

func CreateAuthContext(c *gin.Context, authService auth.Service, userService userService) (auth.PermissionContext, error) {
	currentUser, err := authService.GetCurrentUser(c)
	if err != nil {
		return auth.PermissionContext{}, constants.ErrorNotFound
	}
	currentUserPermissions, err := userService.GetPermissionList(currentUser.ID)
	if err != nil {
		return auth.PermissionContext{}, constants.ErrorNotFound
	}

	return auth.PermissionContext{
		Permissions: currentUserPermissions,
		UserID:      &currentUser.ID,
	}, nil
}

// CreateAdminAuthContext creates an auth context for an admin user
func CreateAdminAuthContext() auth.PermissionContext {
	return auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}
}
