package test

import (
	"testing"
	"time"

	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/test_support"

	"github.com/LydiaTrack/ground/internal/repository"
	"github.com/LydiaTrack/ground/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userService     service.UserService
	initializedUser = false
)

func initializeUserService() {
	if !initializedUser {
		test_support.TestWithMongo()
		repo := repository.GetUserMongoRepository(repository.GetRoleMongoRepository())

		roleService := service.NewRoleService(repository.GetRoleMongoRepository())

		// Create a new user service instance
		userService = *service.NewUserService(repo, *roleService)
		initializedUser = true
	}
}

func TestUserService(t *testing.T) {
	test_support.TestWithMongo()
	initializeUserService()

	t.Run("Create", testCreateUser)
	t.Run("AddRole", testAddRoleToUser)
	t.Run("RemoveRole", testRemoveRoleFromUser)
	t.Run("CreateAndVerifyUser", testCreateAndVerifyUser)
	t.Run("CreateAndDeleteUser", testCreateAndDeleteUser)
}

func testCreateUser(t *testing.T) {

	birthDate := primitive.NewDateTimeFromTime(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	// Create a new userModel
	command := user.CreateUserCommand{
		Username: "test-create-user-001",
		Password: "test123",
		PersonInfo: &user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
			BirthDate: birthDate,
		},
		ContactInfo: user.ContactInfo{
			Email: "test-create-user-001@example.com",
			PhoneNumber: &user.PhoneNumber{
				AreaCode:    "532",
				Number:      "5232323",
				CountryCode: "+90",
			},
		},
	}
	userModel, err := userService.Create(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating userModel test: %v", err)
	} else {

		if userModel.Username != "test-create-user-001" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.PersonInfo.FirstName != "TestName" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.PersonInfo.LastName != "Test Lastname" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.PersonInfo.BirthDate != birthDate {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.ContactInfo.Email != "test-create-user-001@example.com" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.ContactInfo.PhoneNumber.AreaCode != "532" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.ContactInfo.PhoneNumber.Number != "5232323" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.ContactInfo.PhoneNumber.CountryCode != "+90" {
			t.Errorf("Error creating userModel: %v", err)
		}
	}

	// Check user is exists
	exists, err := userService.Exists(userModel.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error checking user exists: %v", err)
	}

	if !exists {
		t.Errorf("Error checking user exists: %v", err)
	}

	// Check user is exists by username
	existsByUsername, err := userService.ExistsByUsername("test-create-user-001", auth.CreateAdminAuthContext())
	if err != nil {
		t.Errorf("Error checking user exists by username: %v", err)
	}

	if !existsByUsername {
		t.Errorf("Error checking user exists by username: %v", err)
	}

	// Check user exists by email address
	existsByEmail, err := userService.ExistsByEmail("test-create-user-001@example.com", auth.CreateAdminAuthContext())
	if err != nil {
		t.Errorf("Error checking user exists by email: %v", err)
	}

	if !existsByEmail {
		t.Errorf("Error checking user exists by email: %v", err)
	}

	// Get user
	userModel, err = userService.Get(userModel.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error getting user: %v", err)
	}

	if userModel.Username != "test-create-user-001" {
		t.Errorf("Error getting user: %v", err)
	}
}

func testAddRoleToUser(t *testing.T) {
	// Create a new userModel
	command := user.CreateUserCommand{
		Username: "test-add-role-to-user-001",
		Password: "test123",
		PersonInfo: &user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
			BirthDate: primitive.NewDateTimeFromTime(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		ContactInfo: user.ContactInfo{
			Email:       "test-add-role-to-user-001@example.com",
			PhoneNumber: nil,
		},
	}

	userModel, err := userService.Create(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating userModel: %v", err)
	}

	roleService := service.NewRoleService(repository.GetRoleMongoRepository())

	// Create a new role
	roleCommand := role.CreateRoleCommand{
		Name: "test-add-role-to-user-001",
		Tags: []string{"testTag"},
		Info: "Test Tag Create",
		Permissions: []auth.Permission{
			auth.AdminPermission,
		},
	}

	roleModel, err := roleService.Create(roleCommand, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating role: %v", err)
	}

	// Check if the role is created
	roleModelAfterCreate, err := roleService.Get(roleModel.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error getting role: %v", roleModelAfterCreate)
	}

	// Add role to user
	addRoleToUserCmd := user.AddRoleToUserCommand{
		UserID: userModel.ID,
		RoleID: roleModelAfterCreate.ID,
	}

	err = userService.AddRole(addRoleToUserCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error adding role to user: %v", err)
	}

	// Check if the role is added to the user
	result, err := userService.GetRoles(userModel.ID, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error getting user roles: %v", err)
	}
	roles := result.Data

	if len(roles) == 0 {
		t.Errorf("Error getting user roles: %v", err)
	}

	if roles[0].ID != roleModel.ID {
		t.Errorf("Error getting user roles: %v", err)
	}

	// Get user's permissions
	permissions, err := userService.GetPermissionList(userModel.ID)

	if err != nil {
		t.Errorf("Error getting user permissions: %v", err)
	}

	if len(permissions) == 0 {
		t.Errorf("Error getting user permissions: %v", err)
	}
}

func testRemoveRoleFromUser(t *testing.T) {
	// Create a new userModel
	command := user.CreateUserCommand{
		Username: "test-remove-role-from-user-001",
		Password: "test123",
		PersonInfo: &user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
		},
		ContactInfo: user.ContactInfo{
			Email:       "test-remove-role-from-user-001@example.com",
			PhoneNumber: nil,
		},
	}

	userModel, err := userService.Create(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating userModel: %v", err)
	}

	roleService := service.NewRoleService(repository.GetRoleMongoRepository())

	// Create a new role
	roleCommand := role.CreateRoleCommand{
		Name: "test-remove-role-from-user-001",
		Tags: []string{"testTag"},
		Info: "Test Tag Create",
		Permissions: []auth.Permission{
			{
				Domain: "USER",
				Action: "READ",
			},
		},
	}

	roleModel, err := roleService.Create(roleCommand, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating role: %v", err)
	}

	// Add role to user
	addRoleToUserCmd := user.AddRoleToUserCommand{
		UserID: userModel.ID,
		RoleID: roleModel.ID,
	}

	err = userService.AddRole(addRoleToUserCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error adding role to user: %v", err)
	}

	// Remove role from user
	removeRoleFromUserCmd := user.RemoveRoleFromUserCommand{
		UserID: userModel.ID,
		RoleID: roleModel.ID,
	}

	err = userService.RemoveRole(removeRoleFromUserCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error removing role from user: %v", err)
	}

	// Check if the role is removed from the user
	result, err := userService.GetRoles(userModel.ID, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error getting user roles: %v", err)
	}
	roles := result.Data

	if len(roles) != 0 {
		t.Errorf("Error getting user roles: %v", err)
	}

	// Get user's permissions
	permissions, err := userService.GetPermissionList(userModel.ID)

	if err != nil {
		t.Errorf("Error getting user permissions: %v", err)
	}

	if len(permissions) != 0 {
		t.Errorf("Error getting user permissions: %v", err)
	}
}

func testCreateAndVerifyUser(t *testing.T) {
	// Create a new userModel
	command := user.CreateUserCommand{
		Username: "test-create-verify-user-001",
		Password: "test123",
		PersonInfo: &user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
		},
		ContactInfo: user.ContactInfo{
			Email:       "test-create-verify-user@example.com",
			PhoneNumber: nil,
		},
	}

	_, err := userService.Create(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating userModel: %v", err)
	}

	// Verify user
	_, err = userService.VerifyUser("test-create-verify-user-001", "test123", auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error verifying userModel: %v", err)
	}

	// Verify user with wrong password
	_, err = userService.VerifyUser("test-create-verify-user-001", "test1234", auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err == nil {
		t.Errorf("Error verifying userModel: %v", err)
	}
}

func testCreateAndDeleteUser(t *testing.T) {
	// Create a new userModel
	command := user.CreateUserCommand{
		Username: "test-create-delete-user-001",
		Password: "test123",
		PersonInfo: &user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
		},
		ContactInfo: user.ContactInfo{
			Email:       "test-create-delete-user-001@example.com",
			PhoneNumber: nil,
		},
	}

	userModel, err := userService.Create(command, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating userModel: %v", err)
	}

	// Delete user
	deleteUserCmd := user.DeleteUserCommand{
		ID: userModel.ID,
	}

	err = userService.Delete(deleteUserCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error deleting userModel: %v", err)
	}
}
