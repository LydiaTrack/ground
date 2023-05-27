package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lydia-track-base/internal/domain"
	"lydia-track-base/internal/domain/commands"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/test_support"
	"testing"
	"time"
)

// TestNewUserService Create a new user service instance with UserMongoRepository
func TestNewUserService(t *testing.T) {
	// Initialize mongo test support
	test_support.TestWithMongo()

	repo := repository.NewUserMongoRepository()

	// Create a new user service instance
	NewUserService(repo)
}

// TestCreateUser Create a new user
func TestCreateUser(t *testing.T) {
	// Initialize mongo test support
	test_support.TestWithMongo()

	repo := repository.NewUserMongoRepository()

	// Create a new user service instance
	userService := NewUserService(repo)

	birthDate := primitive.NewDateTimeFromTime(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	// Create a new user
	command := commands.CreateUserCommand{
		Username: "test",
		Password: "test123",
		PersonInfo: domain.PersonInfo{
			FirstName: "TestName",
			LastName:  "Test Lastname",
			Email:     "exampletest@example.com",
			Address:   "Test Address",
			BirthDate: birthDate,
			PhoneNumber: domain.PhoneNumber{
				AreaCode:    "500",
				Number:      "5005050",
				CountryCode: "+90",
			},
		},
	}
	user, err := userService.CreateUser(command)

	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	if user.Username != "test" {
		t.Errorf("Error creating user: %v", err)
	}

	if user.Password != "test123" {
		t.Errorf("Error creating user: %v", err)
	}

	if user.PersonInfo.FirstName != "TestName" {
		t.Errorf("Error creating user: %v", err)
	}

	if user.PersonInfo.LastName != "Test Lastname" {
		t.Errorf("Error creating user: %v", err)
	}

	if user.PersonInfo.Email != "exampletest@example.com" {
		t.Errorf("Error creating user: %v", err)
	}

	if user.PersonInfo.Address != "Test Address" {
		t.Errorf("Error creating user: %v", err)
	}

	if user.PersonInfo.BirthDate != birthDate {
		t.Errorf("Error creating user: %v", err)
	}

	if user.PersonInfo.PhoneNumber.AreaCode != "500" {
		t.Errorf("Error creating user: %v", err)
	}

	if user.PersonInfo.PhoneNumber.Number != "5005050" {
		t.Errorf("Error creating user: %v", err)
	}

	if user.PersonInfo.PhoneNumber.CountryCode != "+90" {
		t.Errorf("Error creating user: %v", err)
	}
}
