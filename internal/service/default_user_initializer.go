package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lydia-track-base/internal/domain/auth"
	"lydia-track-base/internal/domain/user"
	"lydia-track-base/internal/domain/user/commands"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/utils"
	"os"
	"time"
)

// InitializeDefaultUser initializes the default user with default credentials
func InitializeDefaultUser() error {
	userCreateCmd := commands.CreateUserCommand{
		Username: os.Getenv("DEFAULT_USER_USERNAME"),
		Password: os.Getenv("DEFAULT_USER_PASSWORD"),
		PersonInfo: user.PersonInfo{
			FirstName: "Lydia",
			LastName:  "Admin",
			BirthDate: primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err := NewUserService(repository.GetUserRepository()).
		CreateUser(userCreateCmd, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return err
	}

	utils.Log("Default user created successfully")

	return nil
}
