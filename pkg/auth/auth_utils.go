package auth

import (
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService interface {
	GetPermissionList(userID primitive.ObjectID) ([]Permission, error)
}

func CreateAuthContext(c *gin.Context, authService Service, userService userService) (PermissionContext, error) {
	currentUser, err := authService.GetCurrentUser(c)
	if err != nil {
		return PermissionContext{}, constants.ErrorNotFound
	}
	currentUserPermissions, err := userService.GetPermissionList(currentUser.ID)
	if err != nil {
		return PermissionContext{}, constants.ErrorNotFound
	}

	return PermissionContext{
		Permissions: currentUserPermissions,
		UserID:      &currentUser.ID,
	}, nil
}

// CreateAdminAuthContext creates an auth context for an admin user
func CreateAdminAuthContext() PermissionContext {
	return PermissionContext{
		Permissions: []Permission{AdminPermission},
		UserID:      nil,
	}
}
