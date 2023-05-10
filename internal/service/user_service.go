package service

import (
	"gopkg.in/mgo.v2/bson"
	"lydia-track-base/internal/domain"
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
func (s UserService) CreateUser(user domain.User) (domain.User, error) {
	// TODO: These kind of operations must be done with specific requests, not by User model itself
	user, err := s.userRepository.SaveUser(user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s UserService) GetUser(id string) (domain.User, error) {
	user, err := s.userRepository.GetUser(bson.ObjectIdHex(id))
	if err != nil {
		return domain.User{}, err
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

func (s UserService) DeleteUser(id string) error {
	err := s.userRepository.DeleteUser(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}
	return nil
}

type UserRepository interface {
	// SaveUser saves a user
	SaveUser(user domain.User) (domain.User, error)
	// GetUser gets a user by id
	GetUser(id bson.ObjectId) (domain.User, error)
	// ExistsUser checks if a user exists
	ExistsUser(id bson.ObjectId) (bool, error)
	// DeleteUser deletes a user by id
	DeleteUser(id bson.ObjectId) error
}
