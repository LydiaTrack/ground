package test

import (
	"testing"

	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/domain/role"
	"github.com/LydiaTrack/lydia-base/pkg/test_support"

	"github.com/LydiaTrack/lydia-base/internal/permissions"
	"github.com/LydiaTrack/lydia-base/internal/repository"
	"github.com/LydiaTrack/lydia-base/internal/service"
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
		roleService = *service.NewRoleService(repo)
		initializedRole = true
	}
}

func TestRoleService(t *testing.T) {
	test_support.TestWithMongo()
	initializeRoleService()

	t.Run("CreateRole", testCreateRole)
	t.Run("CannotCreateRoleWithSameName", testCannotCreateRoleWithSameName)
	t.Run("DeleteRole", testDeleteRole)
}

func testCreateRole(t *testing.T) {

	command := role.CreateRoleCommand{
		Name: "testCreate123",
		Tags: []string{"testTag"},
		Info: "Test Tag Create",
	}

	role, err := roleService.CreateRole(command, auth.PermissionContext{
		Permissions: []auth.Permission{permissions.RoleCreatePermission},
		UserId:      nil,
	})

	if err != nil {
		t.Errorf("Error creating role: %s", err)
	}

	if role.Name != command.Name {
		t.Errorf("Expected role name: %s, got: %s", command.Name, role.Name)
	}

	// Check if the role is created or not by existence control
	exists, err := roleService.ExistsRole(role.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})

	if err != nil {
		t.Errorf("Error checking role: %s", err)
	}

	if !exists {
		t.Errorf("Expected role not exists")
	}

	// Check if the role is created or not by getting the role
	role, err = roleService.GetRole(role.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})

	if err != nil {
		t.Errorf("Error getting role: %s", err)
	}

	if role.Name != command.Name {
		t.Errorf("Expected role name: %s, got: %s", command.Name, role.Name)
	}
}

func testCannotCreateRoleWithSameName(t *testing.T) {
	command := role.CreateRoleCommand{
		Name: "testCreate123",
		Tags: []string{"testTag"},
		Info: "Test Tag Create",
	}

	_, err := roleService.CreateRole(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})

	if err == nil {
		t.Errorf("Expected error creating role")
	}

	// Create a new role with the same name
	command = role.CreateRoleCommand{
		Name: "testCreate123",
		Tags: []string{"testTag1", "testTag2"},
		Info: "Test Tag Create123",
	}

	_, err = roleService.CreateRole(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})

	if err == nil {
		t.Errorf("Expected error creating role")
	}
}

func testDeleteRole(t *testing.T) {
	command := role.CreateRoleCommand{
		Name: "testDelete",
		Tags: []string{"testTag"},
		Info: "Test Tag Delete",
	}

	role, err := roleService.CreateRole(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})

	if err != nil {
		t.Errorf("Error creating role: %s", err)
	}

	err = roleService.DeleteRole(role.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})

	if err != nil {
		t.Errorf("Error deleting role: %s", err)
	}

	// Check if the role is deleted or not by existence control
	exists, err := roleService.ExistsRole(role.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})

	if exists {
		t.Errorf("Expected role exists")
	}
}
