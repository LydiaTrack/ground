package initializers

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/internal/domain/role"
	"github.com/LydiaTrack/lydia-base/internal/repository"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/LydiaTrack/lydia-base/internal/utils"
	"os"
)

func InitializeDefaultRole() error {

	// While using remote connection for MongoDB instead of container, the role can be exist in the database.
	// In this case, the default role will not be created.
	isExists := repository.GetRoleRepository().ExistsByRolename(os.Getenv("DEFAULT_ROLE_NAME"))
	if isExists {
		utils.Log("Default role already exists")
		return nil
	}
	roleCreateCmd := role.CreateRoleCommand{
		Name: os.Getenv("DEFAULT_ROLE_NAME"),
		Tags: []string{os.Getenv("DEFAULT_ROLE_TAG")},
		Info: os.Getenv("DEFAULT_ROLE_INFO"),
	}

	_, err := service.NewRoleService(repository.GetRoleRepository()).CreateRole(roleCreateCmd, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return err
	}

	utils.Log("Default role created successfully")

	return nil
}
