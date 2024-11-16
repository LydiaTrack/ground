package initializers

import (
	"os"
	"time"

	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/domain/user"

	"github.com/LydiaTrack/ground/internal/log"
	"github.com/LydiaTrack/ground/internal/repository"
	"github.com/LydiaTrack/ground/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InitializeDefaultUser initializes the default user with default credentials
func InitializeDefaultUser() error {
	// While using remote connection for MongoDB instead of container, the user can be exist in the database.
	// In this case, the default user will not be created.
	roleService := service.NewRoleService(repository.GetRoleRepository())
	userService := service.NewUserService(repository.GetUserRepository(), *roleService)
	isExists := userService.ExistsByUsername(os.Getenv("DEFAULT_USER_USERNAME"))

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
			ContactInfo: user.ContactInfo{
				Email:       "lydia@lydiaadmin.com",
				PhoneNumber: nil,
			},
		}

		createdUser, err := userService.CreateUser(userCreateCmd, auth.PermissionContext{
			Permissions: []auth.Permission{auth.AdminPermission},
			UserId:      nil,
		})
		if err != nil {
			return err
		}

		// Add admin roles to the default user
		err = addAdminRolesToUser(createdUser, *roleService, *userService)
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
		for _, roleId := range *userRoleIds {
			roleModel, err := roleService.GetRole(roleId.Hex(), auth.PermissionContext{
				Permissions: []auth.Permission{auth.AdminPermission},
				UserId:      nil,
			})
			if err != nil {
				return err
			}
			permissions := roleModel.Permissions
			for _, permission := range permissions {
				if permission == auth.AdminPermission {
					hasAdminRoles = true
					break
				}
			}
		}
	}

	// After checking if the user has admin roles, return if the user has admin roles
	if hasAdminRoles {
		return nil
	}

	// if admin user does not have admin roles, add admin roles

	// Check if admin role exists
	existsRole := roleService.ExistsByName("ADMIN", auth.PermissionContext{
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
	} else {
		// If admin role exists, get the role
		model, err := roleService.GetRoleByName("ADMIN", auth.PermissionContext{
			Permissions: []auth.Permission{auth.AdminPermission},
			UserId:      nil,
		})
		if err != nil {
			return err
		}
		roleModel = model

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
		Name:        "ADMIN",
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
