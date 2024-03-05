package initializers

import (
	"github.com/Lydia/lydia-base/internal/domain/auth"
	"github.com/Lydia/lydia-base/internal/domain/role"
	"github.com/Lydia/lydia-base/internal/repository"
	"github.com/Lydia/lydia-base/internal/service"
	"github.com/Lydia/lydia-base/internal/utils"
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
