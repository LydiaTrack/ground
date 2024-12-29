package test

import (
	"testing"

	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/test_support"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LydiaTrack/ground/internal/permissions"
	"github.com/LydiaTrack/ground/internal/repository"
	"github.com/LydiaTrack/ground/internal/service"
)

var (
	roleService     service.RoleService
	initializedRole = false
)

func initializeRoleService() {
	if !initializedRole {
		test_support.TestWithMongo()
		repo := repository.GetRoleMongoRepository()

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
		Name: "testCreateRole",
		Tags: []string{"testTag"},
		Info: "Test Tag Create",
	}

	roleModel, err := roleService.CreateRole(command, auth.PermissionContext{
		Permissions: []auth.Permission{permissions.RoleCreatePermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating roleModel: %s", err)
	}

	if roleModel.Name != command.Name {
		t.Errorf("Expected roleModel name: %s, got: %s", command.Name, roleModel.Name)
	}

	// Check if the roleModel is created or not by existence control
	exists, err := roleService.Exists(roleModel.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error checking roleModel: %s", err)
	}

	if !exists {
		t.Errorf("Expected roleModel not exists")
	}

	// Check if the roleModel is created or not by getting the roleModel
	roleModel, err = roleService.GetRole(roleModel.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error getting roleModel: %s", err)
	}

	if roleModel.Name != command.Name {
		t.Errorf("Expected roleModel name: %s, got: %s", command.Name, roleModel.Name)
	}
}

func testCannotCreateRoleWithSameName(t *testing.T) {
	command := role.CreateRoleCommand{
		Name: "testCannotCreateRole",
		Tags: []string{"testTag"},
		Info: "Test Tag Create",
	}

	_, err := roleService.CreateRole(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating roleModel: %s", err)
	}

	// Create a new role with the same name
	command = role.CreateRoleCommand{
		Name: "testCannotCreateRole",
		Tags: []string{"testTag1", "testTag2"},
		Info: "Test Tag Create123",
	}

	_, err = roleService.CreateRole(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
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

	roleModel, err := roleService.CreateRole(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating roleModel: %s", err)
	}

	err = roleService.DeleteRole(roleModel.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error deleting roleModel: %s", err)
	}

	// Check if the roleModel is deleted or not by existence control
	exists, err := roleService.Exists(roleModel.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	// If error is not nil, and error is not mongo.ErrNoDocuments, then it is an unexpected error
	if err != nil && err != mongo.ErrNoDocuments {
		t.Errorf("Error checking roleModel: %s", err)
	}

	if exists {
		t.Errorf("Expected roleModel exists")
	}
}
