package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lydia-track-base/internal/domain"
	"lydia-track-base/internal/domain/commands"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
	"time"
)

// InitializeDefaultUser initializes the default user with default credentials
func InitializeDefaultUser() error {
	userCreateCmd := commands.CreateUserCommand{
		Username: "lydia",
		Password: "P@ssw0rd",
		PersonInfo: domain.PersonInfo{
			FirstName: "Lydia",
			LastName:  "Admin",
			BirthDate: primitive.
				DateTime(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond)),
		},
	}

	_, err := service.
		NewUserService(repository.GetUserRepository()).
		CreateUser(userCreateCmd)
	if err != nil {
		return err
	}

	return nil
}
