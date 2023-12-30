package test

import (
	"lydia-track-base/internal/domain/auth"
	"lydia-track-base/internal/domain/user"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
	"lydia-track-base/internal/test_support"
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
		userService = service.NewUserService(repo)
		initializedUser = true
	}
}

//func TestMain(m *testing.M) {
//	initializeUserService()
//	m.Run()
//}

// TestNewUserService Create a new user service instance with UserMongoRepository
func TestNewUserService(t *testing.T) {
	test_support.TestWithMongo()

	// Check for user service is initializedUser or not
	if !initializedUser {
		t.Errorf("Error initializing user service")
	}
}

// TestCreateUser Create a new user
func TestCreateUser(t *testing.T) {
	test_support.TestWithMongo()

	// Create a new userModel service instance
	birthDate := primitive.NewDateTimeFromTime(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	// Create a new userModel
	command := user.CreateUserCommand{
		Username: "test",
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

		if userModel.Username != "test" {
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
}
