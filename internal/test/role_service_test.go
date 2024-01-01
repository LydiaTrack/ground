package test

import (
	"lydia-track-base/internal/domain/auth"
	"lydia-track-base/internal/domain/role"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
	"lydia-track-base/internal/test_support"
	"testing"
)

var (
	roleService     service.RoleService
	initializedRole = false
)

func initializeRoleService() {
	if !initializedRole {
		test_support.TestWithMongo()
		repo := repository.GetRoleRepository()

		// Create a new role service instance
		roleService = service.NewRoleService(repo)
		initializedRole = true
	}
}

func TestRoleService(t *testing.T) {
	test_support.TestWithMongo()
	initializeRoleService()

	t.Run("CreateRole", testCreateRole)
	t.Run("DeleteRole", testDeleteRole)
}

func testCreateRole(t *testing.T) {
	t.Run("CreateRole", func(t *testing.T) {
		test_support.TestWithMongo()
		initializeRoleService()

		command := role.CreateRoleCommand{
			Name: "testCreate123",
			Tags: []string{"testTag"},
			Info: "Test Tag Create",
		}

		role, err := roleService.CreateRole(command, []auth.Permission{auth.AdminPermission})

		if err != nil {
			t.Errorf("Error creating role: %s", err)
		}

		if role.Name != command.Name {
			t.Errorf("Expected role name: %s, got: %s", command.Name, role.Name)
		}

		// Check if the role is created or not by existence control
		exists, err := roleService.ExistsRole(role.ID.Hex(), []auth.Permission{auth.AdminPermission})

		if err != nil {
			t.Errorf("Error checking role: %s", err)
		}

		if !exists {
			t.Errorf("Expected role not exists")
		}
	})
}

func testDeleteRole(t *testing.T) {
	t.Run("DeleteRole", func(t *testing.T) {
		test_support.TestWithMongo()
		initializeRoleService()

		command := role.CreateRoleCommand{
			Name: "testDelete",
			Tags: []string{"testTag"},
			Info: "Test Tag Delete",
		}

		role, err := roleService.CreateRole(command, []auth.Permission{auth.AdminPermission})

		if err != nil {
			t.Errorf("Error creating role: %s", err)
		}

		err = roleService.DeleteRole(role.ID.Hex(), []auth.Permission{auth.AdminPermission})

		if err != nil {
			t.Errorf("Error deleting role: %s", err)
		}

		// Check if the role is deleted or not by existence control
		exists, err := roleService.ExistsRole(role.ID.Hex(), []auth.Permission{auth.AdminPermission})

		if exists {
			t.Errorf("Expected role exists")
		}
	})
}
