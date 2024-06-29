package initializers

import (
	"os"
	"time"

	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/domain/role"
	"github.com/LydiaTrack/lydia-base/pkg/domain/user"

	"github.com/LydiaTrack/lydia-base/internal/log"
	"github.com/LydiaTrack/lydia-base/internal/repository"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InitializeDefaultUser initializes the default user with default credentials
func InitializeDefaultUser() error {
	// While using remote connection for MongoDB instead of container, the user can be exist in the database.
	// In this case, the default user will not be created.
	roleService := service.NewRoleService(repository.GetRoleRepository())
	userService := service.NewUserService(repository.GetUserRepository(), *roleService)
	isExists, err := userService.ExistsByUsername(os.Getenv("DEFAULT_USER_USERNAME"))
	if err != nil {
		return err
	}

	if isExists {
		log.Log("Default user already exists")
		userModel, err := repository.GetUserRepository().GetUserByUsername(os.Getenv("DEFAULT_USER_USERNAME"))
		if err != nil {
			return err
		}

		// Add admin roles to the default user
		err = addAdminRolesToUser(userModel, *roleService, *userService)
		if err != nil {
			return err
		}

	} else {
		// If default user does not exist, create default user
		userCreateCmd := user.CreateUserCommand{
			Username: os.Getenv("DEFAULT_USER_USERNAME"),
			Password: os.Getenv("DEFAULT_USER_PASSWORD"),
			PersonInfo: &user.PersonInfo{
				FirstName: "Lydia",
				LastName:  "Admin",
				BirthDate: primitive.NewDateTimeFromTime(time.Now()),
			},
		}

		userCreateResponse, err := userService.CreateUser(userCreateCmd, auth.PermissionContext{
			Permissions: []auth.Permission{auth.AdminPermission},
			UserId:      nil,
		})
		if err != nil {
			return err
		}

		userModel := user.Model{
			ID:          userCreateResponse.ID,
			Username:    userCreateResponse.Username,
			Password:    "",
			PersonInfo:  &user.PersonInfo{},
			CreatedDate: time.Time{},
			Version:     0,
			RoleIds:     nil,
		}

		// Add admin roles to the default user
		err = addAdminRolesToUser(userModel, *roleService, *userService)
		if err != nil {
			return err
		}
	}

	log.Log("Default user created successfully")

	return nil
}

func addAdminRolesToUser(userModel user.Model, roleService service.RoleService, userService service.UserService) error {
	// Check if userModel has admin role
	userRoleIds := userModel.RoleIds
	hasAdminRoles := false
	if userRoleIds == nil {
		log.Log("Default user does not have any roles, adding admin roles...")
	} else {
		for _, roleId := range userRoleIds {
			roleModel, err := roleService.GetRole(roleId.Hex(), auth.PermissionContext{
				Permissions: []auth.Permission{auth.AdminPermission},
				UserId:      nil,
			})
			if err != nil {
				return err
			}
			if roleModel.Name == "LYDIA_ADMIN" {
				log.Log("Default user has Lydia Admin role")
				hasAdminRoles = true
				break
			}
		}
	}

	if hasAdminRoles {
		return nil
	}

	// if admin user does not have admin roles, add admin roles

	// Check if admin role exists
	existsRole := roleService.ExistsByName("LYDIA_ADMIN", auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})
	roleModel := role.Model{}
	if !existsRole {
		// If admin role does not exist, create admin role
		rm, err := createAdminRole(roleService)
		if err != nil {
			return err
		}
		roleModel = rm
	}

	addRoleCmd := user.AddRoleToUserCommand{
		UserID: userModel.ID,
		RoleID: roleModel.ID,
	}
	err := userService.AddRoleToUser(addRoleCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})
	if err != nil {
		return err
	}
	return nil
}

// createAdminRole creates the admin role for the default user
func createAdminRole(roleService service.RoleService) (role.Model, error) {
	createRoleCmd := role.CreateRoleCommand{
		Name:        "LYDIA_ADMIN",
		Tags:        nil,
		Info:        "Lydia Admin role for the default user",
		Permissions: []auth.Permission{auth.AdminPermission},
	}
	roleModel, err := roleService.CreateRole(createRoleCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}
