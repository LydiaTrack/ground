package initializers

import (
	"lydia-track-base/internal/domain/auth"
	"lydia-track-base/internal/domain/role"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
	"lydia-track-base/internal/utils"
	"os"
)

func InitializeDefaultRole() error {
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
