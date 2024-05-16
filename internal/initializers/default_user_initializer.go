package initializers

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/internal/domain/role"
	"github.com/LydiaTrack/lydia-base/internal/domain/user"
	"github.com/LydiaTrack/lydia-base/internal/repository"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/LydiaTrack/lydia-base/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

// InitializeDefaultUser initializes the default user with default credentials
func InitializeDefaultUser() error {
	// While using remote connection for MongoDB instead of container, the user can be exist in the database.
	// In this case, the default user will not be created.
	userService := service.NewUserService(repository.GetUserRepository())
	roleService := service.NewRoleService(repository.GetRoleRepository())
	isExists, err := userService.ExistsByUsername(os.Getenv("DEFAULT_USER_USERNAME"), []auth.Permission{auth.AdminPermission})
	if err != nil {
		return err
	}

	if isExists {
		utils.Log("Default user already exists")
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
			PersonInfo: user.PersonInfo{
				FirstName: "Lydia",
				LastName:  "Admin",
				BirthDate: primitive.NewDateTimeFromTime(time.Now()),
			},
		}

		userCreateResponse, err := userService.CreateUser(userCreateCmd, []auth.Permission{auth.AdminPermission})
		if err != nil {
			return err
		}

		userModel := user.Model{
			ID:          userCreateResponse.ID,
			Username:    userCreateResponse.Username,
			Password:    "",
			PersonInfo:  user.PersonInfo{},
			CreatedDate: time.Time{},
			Version:     0,
			RoleIds:     nil,
		}

		// Add admin roles to the default user
		err = addAdminRolesToUser(userModel, *roleService, *userService)
	}

	utils.Log("Default user created successfully")

	return nil
}

func addAdminRolesToUser(userModel user.Model, roleService service.RoleService, userService service.UserService) error {
	// Check if userModel has admin role
	userRoleIds := userModel.RoleIds
	hasAdminRoles := false
	if userRoleIds == nil || len(userRoleIds) == 0 {
		utils.Log("Default user does not have any roles, adding admin roles...")
	} else {
		for _, roleId := range userRoleIds {
			roleModel, err := roleService.GetRole(roleId.Hex(), []auth.Permission{auth.AdminPermission})
			if err != nil {
				return err
			}
			if roleModel.Name == "LYDIA_ADMIN" {
				utils.Log("Default user has Lydia Admin role")
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
	existsRole := roleService.ExistsByRolename("LYDIA_ADMIN", []auth.Permission{auth.AdminPermission})
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
	err := userService.AddRoleToUser(addRoleCmd, []auth.Permission{auth.AdminPermission})
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
	roleModel, err := roleService.CreateRole(createRoleCmd, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return role.Model{}, err
	}
	return roleModel, nil
}
