package auth

import (
	"fmt"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/gin-gonic/gin"
	"time"
)

type userService interface {
	GetPermissionList(userModel user.Model) ([]Permission, error)
}

func CreateAuthContext(c *gin.Context, authService Service, userService userService) (PermissionContext, error) {
	now := time.Now()
	currentUser, err := authService.GetCurrentUser(c)
	if err != nil {
		return PermissionContext{}, constants.ErrorNotFound
	}
	elapsedCurrentUser := time.Since(now)
	currentUserPermissions, err := userService.GetPermissionList(currentUser)
	if err != nil {
		return PermissionContext{}, constants.ErrorNotFound
	}
	elapsedPermissions := time.Since(now) - elapsedCurrentUser
	elapsed := time.Since(now)
	if elapsed > 100*time.Millisecond {
		fmt.Printf("Auth context took too long to create, elapsed %v, user: %v, permissions: %v\n", elapsed, elapsedCurrentUser, elapsedPermissions)
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
