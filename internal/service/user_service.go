package service

import (
	"errors"
	"github.com/LydiaTrack/lydia-base/internal/domain/auth"
	"github.com/LydiaTrack/lydia-base/internal/domain/role"
	"github.com/LydiaTrack/lydia-base/internal/domain/user"
	"github.com/LydiaTrack/lydia-base/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
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

type UserRepository interface {
	// SaveUser saves a user
	SaveUser(user user.Model) (user.Model, error)
	// GetUser gets a user by id
	GetUser(id bson.ObjectId) (user.Model, error)
	// GetUserByUsername gets a user by username
	GetUserByUsername(username string) (user.Model, error)
	// ExistsUser checks if a user exists
	ExistsUser(id bson.ObjectId) (bool, error)
	// DeleteUser deletes a user by id
	DeleteUser(id bson.ObjectId) error
	// ExistsByUsername checks if a user exists by username
	ExistsByUsername(username string) bool
	// AddRoleToUser adds a role to a user
	AddRoleToUser(userID bson.ObjectId, roleID bson.ObjectId) error
	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(userID bson.ObjectId, roleID bson.ObjectId) error
	// GetUserRoles gets the roles of a user
	GetUserRoles(userID bson.ObjectId) ([]role.Model, error)
}

func (s UserService) CreateUser(command user.CreateUserCommand, permissions []auth.Permission) (user.CreateResponse, error) {
	if !CheckPermission(permissions, user.CreatePermission) {
		return user.CreateResponse{}, errors.New("not permitted")
	}

	// Validate user
	// Map command to user
	userModel := user.NewUser(bson.NewObjectId().Hex(), command.Username,
		command.Password, command.PersonInfo, time.Now(), 1)
	if err := userModel.Validate(); err != nil {
		return user.CreateResponse{}, err
	}
	userExists := s.userRepository.ExistsByUsername(userModel.Username)

	if userExists {
		return user.CreateResponse{}, errors.New("user already exists")
	}

	userModel, err := beforeCreateUser(userModel)
	if err != nil {
		return user.CreateResponse{}, err
	}

	savedUser, err := s.userRepository.SaveUser(userModel)
	if err != nil {
		return user.CreateResponse{}, err
	}
	savedUser, err = afterCreateUser(savedUser)
	if err != nil {
		return user.CreateResponse{}, err
	}

	response := user.CreateResponse{
		ID:          savedUser.ID,
		Username:    savedUser.Username,
		PersonInfo:  savedUser.PersonInfo,
		CreatedDate: savedUser.CreatedDate,
		Version:     savedUser.Version,
	}
	utils.Log("User %s created successfully", response.Username)
	return response, nil
}

// beforeCreateUser is a hook that is called before creating a user
func beforeCreateUser(userModel user.Model) (user.Model, error) {
	// Hash user password before saving
	hashedPassword, err := hashPassword(userModel.Password)
	if err != nil {
		return user.Model{}, err
	}

	userModel.Password = hashedPassword
	return userModel, nil
}

// afterCreateUser is a hook that is called after creating a user
func afterCreateUser(user user.Model) (user.Model, error) {
	return user, nil
}

func (s UserService) GetUser(id string, permissions []auth.Permission) (user.Model, error) {
	if !CheckPermission(permissions, user.ReadPermission) {
		return user.Model{}, errors.New("not permitted")
	}

	userModel, err := s.userRepository.GetUser(bson.ObjectIdHex(id))
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

func (s UserService) ExistsUser(id string, permissions []auth.Permission) (bool, error) {
	if !CheckPermission(permissions, user.ReadPermission) {
		return false, errors.New("not permitted")
	}

	exists, err := s.userRepository.ExistsUser(bson.ObjectIdHex(id))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s UserService) DeleteUser(command user.DeleteUserCommand, permissions []auth.Permission) error {
	if !CheckPermission(permissions, user.DeletePermission) {
		return errors.New("not permitted")
	}

	existsUser, err := s.ExistsUser(command.ID.Hex(), permissions)
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

// hashPassword hashes a password using bcrypt
func hashPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyUser verifies a user by username and password
func (s UserService) VerifyUser(username string, password string, permissions []auth.Permission) (user.Model, error) {
	if !CheckPermission(permissions, user.ReadPermission) {
		return user.Model{}, errors.New("not permitted")
	}

	// Get the user by username
	userModel, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return user.Model{}, err
	}

	// Compare the passwords
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(password))
	if err != nil {
		utils.LogError("Error comparing passwords: " + err.Error())
		return user.Model{}, err
	}

	return userModel, nil
}

// ExistsByUsername gets a user by username
func (s UserService) ExistsByUsername(username string, permissions []auth.Permission) (bool, error) {
	if !CheckPermission(permissions, user.ReadPermission) {
		return false, errors.New("not permitted")
	}
	return s.userRepository.ExistsByUsername(username), nil
}

// AddRoleToUser adds a role to a user
func (s UserService) AddRoleToUser(command user.AddRoleToUserCommand, permissions []auth.Permission) error {
	if !CheckPermission(permissions, user.UpdatePermission) {
		return errors.New("not permitted")
	}
	return s.userRepository.AddRoleToUser(command.UserID, command.RoleID)
}

// RemoveRoleFromUser removes a role from a user
func (s UserService) RemoveRoleFromUser(command user.RemoveRoleFromUserCommand, permissions []auth.Permission) error {
	if !CheckPermission(permissions, user.UpdatePermission) {
		return errors.New("not permitted")
	}
	return s.userRepository.RemoveRoleFromUser(command.UserID, command.RoleID)
}

// GetUserRoles gets the roles of a user
func (s UserService) GetUserRoles(userID bson.ObjectId, permissions []auth.Permission) ([]role.Model, error) {
	if !CheckPermission(permissions, user.ReadPermission) {
		return nil, errors.New("not permitted")
	}
	return s.userRepository.GetUserRoles(userID)
}

// GetUserPermissions gets the permissions of a user
func (s UserService) GetUserPermissions(userID bson.ObjectId) ([]auth.Permission, error) {
	userRoles, err := s.GetUserRoles(userID, []auth.Permission{auth.AdminPermission})
	if err != nil {
		return nil, err
	}

	var userPermissions []auth.Permission
	for _, userRole := range userRoles {
		userPermissions = append(userPermissions, userRole.Permissions...)
	}

	return userPermissions, nil
}
