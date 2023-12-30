package test

import (
	"lydia-track-base/internal/domain/auth"
	"lydia-track-base/internal/domain/role/commands"
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

func TestMain(m *testing.M) {
	initializeRoleService()
	m.Run()
}

func TestNewRoleService(t *testing.T) {
	test_support.TestWithMongo()

	if !initializedRole {
		t.Errorf("Error initializing role service")
	}
}

func TestCreateRole(t *testing.T) {
	test_support.TestWithMongo()

	command := commands.CreateRoleCommand{
		Name: "testCreate",
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

	// Check for if role is created or not by exist control
	exists, err := roleService.ExistsRole(role.ID.Hex(), []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error checking role: %s", err)
	}

	if !exists {
		t.Errorf("Expected role not exists")
	}
}

func TestDeleteRole(t *testing.T) {
	test_support.TestWithMongo()

	command := commands.CreateRoleCommand{
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

	// Check for if role is deleted or not by exist control
	exists, err := roleService.ExistsRole(role.ID.Hex(), []auth.Permission{auth.AdminPermission})

	if exists {
		t.Errorf("Expected role exists")
	}
}
