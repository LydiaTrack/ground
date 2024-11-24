package initializers

import (
	"os"

	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/role"

	"github.com/LydiaTrack/ground/internal/log"
	"github.com/LydiaTrack/ground/internal/repository"
	"github.com/LydiaTrack/ground/internal/service"
)

func InitializeDefaultRole() error {

	// While using remote connection for MongoDB instead of container, the role can be exist in the database.
	// In this case, the default role will not be created.
	isExists := repository.GetRoleRepository().ExistsByName(os.Getenv("DEFAULT_ROLE_NAME"))
	if isExists {
		log.Log("Default role already exists")
		return nil
	}
	roleCreateCmd := role.CreateRoleCommand{
		Name: os.Getenv("DEFAULT_ROLE_NAME"),
		Tags: []string{os.Getenv("DEFAULT_ROLE_TAG")},
		Info: os.Getenv("DEFAULT_ROLE_INFO"),
	}

	_, err := service.NewRoleService(repository.GetRoleRepository()).CreateRole(roleCreateCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		return err
	}

	log.Log("Default role created successfully")

	return nil
}
