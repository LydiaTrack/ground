package test

import (
	"github.com/LydiaTrack/lydia-base/auth"
	"github.com/LydiaTrack/lydia-base/internal/domain/role"
	"github.com/LydiaTrack/lydia-base/internal/domain/user"
	"github.com/LydiaTrack/lydia-base/internal/repository"
	"github.com/LydiaTrack/lydia-base/internal/service"
	"github.com/LydiaTrack/lydia-base/test_support"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userService     service.UserService
	initializedUser = false
)

func initializeUserService() {
	if !initializedUser {
		test_support.TestWithMongo()
		repo := repository.GetUserRepository()

		// Create a new user service instance
		userService = *service.NewUserService(repo)
		initializedUser = true
	}
}

func TestUserService(t *testing.T) {
	test_support.TestWithMongo()
	initializeUserService()
	initializeRoleService()

	t.Run("CreateUser", testCreateUser)
	t.Run("AddRoleToUser", testAddRoleToUser)
	t.Run("RemoveRoleFromUser", testRemoveRoleFromUser)
	t.Run("CreateAndVerifyUser", testCreateAndVerifyUser)
	t.Run("CreateAndDeleteUser", testCreateAndDeleteUser)
}

func testCreateUser(t *testing.T) {

	birthDate := primitive.NewDateTimeFromTime(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	// Create a new userModel
	command := user.CreateUserCommand{
		Username: "test-create-user-001",
		Password: "test123",
		PersonInfo: user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
			Email:     "exampletest@example.com",
			Address:   "Test Address",
			BirthDate: birthDate,
			PhoneNumber: user.PhoneNumber{
				AreaCode:    "500",
				Number:      "5005050",
				CountryCode: "+90",
			},
		},
	}
	userModel, err := userService.CreateUser(command, []auth.Permission{auth.AdminPermission})

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

		if userModel.PersonInfo.Email != "exampletest@example.com" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.PersonInfo.Address != "Test Address" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.PersonInfo.BirthDate != birthDate {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.PersonInfo.PhoneNumber.AreaCode != "500" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.PersonInfo.PhoneNumber.Number != "5005050" {
			t.Errorf("Error creating userModel: %v", err)
		}

		if userModel.PersonInfo.PhoneNumber.CountryCode != "+90" {
			t.Errorf("Error creating userModel: %v", err)
		}
	}

	// Check user is exists
	exists, err := userService.ExistsUser(userModel.ID.Hex(), []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error checking user exists: %v", err)
	}

	if !exists {
		t.Errorf("Error checking user exists: %v", err)
	}

	// Check user is exists by username
	existsByUsername, err := userService.ExistsByUsername("test-create-user-001", []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error checking user exists by username: %v", err)
	}

	if !existsByUsername {
		t.Errorf("Error checking user exists by username: %v", err)
	}

	// Get user
	getUserModel, err := userService.GetUser(userModel.ID.Hex(), []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error getting user: %v", err)
	}

	if getUserModel.Username != "test-create-user-001" {
		t.Errorf("Error getting user: %v", err)
	}
}

func testAddRoleToUser(t *testing.T) {
	// Create a new userModel
	command := user.CreateUserCommand{
		Username: "test-add-role-to-user-001",
		Password: "test123",
		PersonInfo: user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
			Email:     "test@example.com",
			Address:   "Test Address",
			BirthDate: primitive.NewDateTimeFromTime(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			PhoneNumber: user.PhoneNumber{
				AreaCode:    "500",
				Number:      "5005050",
				CountryCode: "+90",
			},
		},
	}

	userModel, err := userService.CreateUser(command, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating userModel: %v", err)
	}

	roleService := service.NewRoleService(repository.GetRoleRepository())

	// Create a new role
	roleCommand := role.CreateRoleCommand{
		Name: "test-add-role-to-user-001",
		Tags: []string{"testTag"},
		Info: "Test Tag Create",
		Permissions: []auth.Permission{
			auth.AdminPermission,
		},
	}

	roleModel, err := roleService.CreateRole(roleCommand, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating role: %v", err)
	}

	// Add role to user
	addRoleToUserCmd := user.AddRoleToUserCommand{
		UserID: userModel.ID,
		RoleID: roleModel.ID,
	}

	err = userService.AddRoleToUser(addRoleToUserCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error adding role to user: %v", err)
	}

	// Check if the role is added to the user
	roles, err := userService.GetUserRoles(userModel.ID, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error getting user roles: %v", err)
	}

	if len(roles) == 0 {
		t.Errorf("Error getting user roles: %v", err)
	}

	if roles[0].ID != roleModel.ID {
		t.Errorf("Error getting user roles: %v", err)
	}

	// Get user's permissions
	permissions, err := userService.GetUserPermissionList(userModel.ID)

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
		PersonInfo: user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
			Email:     "test@example.com",
			Address:   "Test Address",
		},
	}

	userModel, err := userService.CreateUser(command, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating userModel: %v", err)
	}

	roleService := service.NewRoleService(repository.GetRoleRepository())

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

	roleModel, err := roleService.CreateRole(roleCommand, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating role: %v", err)
	}

	// Add role to user
	addRoleToUserCmd := user.AddRoleToUserCommand{
		UserID: userModel.ID,
		RoleID: roleModel.ID,
	}

	err = userService.AddRoleToUser(addRoleToUserCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error adding role to user: %v", err)
	}

	// Remove role from user
	removeRoleFromUserCmd := user.RemoveRoleFromUserCommand{
		UserID: userModel.ID,
		RoleID: roleModel.ID,
	}

	err = userService.RemoveRoleFromUser(removeRoleFromUserCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error removing role from user: %v", err)
	}

	// Check if the role is removed from the user
	roles, err := userService.GetUserRoles(userModel.ID, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error getting user roles: %v", err)
	}

	if len(roles) != 0 {
		t.Errorf("Error getting user roles: %v", err)
	}

	// Get user's permissions
	permissions, err := userService.GetUserPermissionList(userModel.ID)

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
		PersonInfo: user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
			Email:     "test@gmail.com",
			Address:   "Test Address",
		},
	}

	_, err := userService.CreateUser(command, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating userModel: %v", err)
	}

	// Verify user
	_, err = userService.VerifyUser("test-create-verify-user-001", "test123", []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error verifying userModel: %v", err)
	}

	// Verify user with wrong password
	_, err = userService.VerifyUser("test-create-verify-user-001", "test1234", []auth.Permission{auth.AdminPermission})

	if err == nil {
		t.Errorf("Error verifying userModel: %v", err)
	}
}

func testCreateAndDeleteUser(t *testing.T) {
	// Create a new userModel
	command := user.CreateUserCommand{
		Username: "test-create-delete-user-001",
		Password: "test123",
		PersonInfo: user.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
			Email:     "test@gmail.com",
			Address:   "Test Address",
		},
	}

	userModel, err := userService.CreateUser(command, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating userModel: %v", err)
	}

	// Delete user
	deleteUserCmd := user.DeleteUserCommand{
		ID: userModel.ID,
	}

	err = userService.DeleteUser(deleteUserCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error deleting userModel: %v", err)
	}
}
