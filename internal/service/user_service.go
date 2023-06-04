package service

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain"
	"lydia-track-base/internal/domain/commands"
	"time"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

// CreateUser TODO: Add permission check
func (s UserService) CreateUser(command commands.CreateUserCommand) (domain.UserModel, error) {
	// TODO: These kind of operations must be done with specific requests, not by UserModel model itself
	// Validate user
	// Map command to user
	user := domain.NewUser(bson.NewObjectId().Hex(), command.Username,
		command.Password, command.PersonInfo, time.Now(), 1)
	if err := user.Validate(); err != nil {
		return user, err
	}

	userExists := s.userRepository.ExistsByUsername(user.Username)

	if userExists {
		return domain.UserModel{}, errors.New("user already exists")
	}

	user, err := s.userRepository.SaveUser(user)
	if err != nil {
		return domain.UserModel{}, err
	}
	return user, nil
}

func (s UserService) GetUser(id string) (domain.UserModel, error) {
	user, err := s.userRepository.GetUser(bson.ObjectIdHex(id))
	if err != nil {
		return domain.UserModel{}, err
	}
	return user, nil
}

func (s UserService) ExistsUser(id string) (bool, error) {
	exists, err := s.userRepository.ExistsUser(bson.ObjectIdHex(id))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s UserService) DeleteUser(command commands.DeleteUserCommand) error {
	existsUser, err := s.ExistsUser(command.ID.Hex())
	if err != nil {
		return err
	}
	if !existsUser {
		return errors.New("user does not exist")
	}

	err = s.userRepository.DeleteUser(command.ID)
	if err != nil {
		return err
	}
	return nil
}

type UserRepository interface {
	// SaveUser saves a user
	SaveUser(user domain.UserModel) (domain.UserModel, error)
	// GetUser gets a user by id
	GetUser(id bson.ObjectId) (domain.UserModel, error)
	// ExistsUser checks if a user exists
	ExistsUser(id bson.ObjectId) (bool, error)
	// DeleteUser deletes a user by id
	DeleteUser(id bson.ObjectId) error
	// ExistsByUsername checks if a user exists by username
	ExistsByUsername(username string) bool
}
