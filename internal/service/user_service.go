package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/LydiaTrack/lydia-base/internal/log"
	"github.com/LydiaTrack/lydia-base/internal/permissions"
	"github.com/LydiaTrack/lydia-base/pkg/auth"
	"github.com/LydiaTrack/lydia-base/pkg/constants"
	"github.com/LydiaTrack/lydia-base/pkg/domain/role"
	"github.com/LydiaTrack/lydia-base/pkg/domain/user"
	"github.com/LydiaTrack/lydia-base/pkg/manager"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository UserRepository
	roleService    RoleService
}

func NewUserService(userRepository UserRepository, roleService RoleService) *UserService {
	return &UserService{
		userRepository: userRepository,
		roleService:    roleService,
	}
}

type UserRepository interface {
	// SaveUser saves a user
	SaveUser(user user.Model) (user.Model, error)
	// GetUsers gets all users
	GetUsers() ([]user.Model, error)
	// GetUser gets a user by id
	GetUser(id primitive.ObjectID) (user.Model, error)
	// GetUserByUsername gets a user by username
	GetUserByUsername(username string) (user.Model, error)
	// ExistsUser checks if a user exists
	ExistsUser(id primitive.ObjectID) (bool, error)
	// DeleteUser deletes a user by id
	DeleteUser(id primitive.ObjectID) error
	// ExistsByUsername checks if a user exists by username
	ExistsByUsernameAndEmail(username string, email string) bool
	// ExistsByUsername checks if a user exists by username
	ExistsByUsername(username string) bool
	// AddRoleToUser adds a role to a user
	AddRoleToUser(userID primitive.ObjectID, roleID primitive.ObjectID) error
	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(userID primitive.ObjectID, roleID primitive.ObjectID) error
	// GetUserRoles gets the roles of a user
	GetUserRoles(userID primitive.ObjectID) ([]role.Model, error)
	// UpdateUser updates a user and returns the updated user
	UpdateUser(id primitive.ObjectID, updateCommand user.UpdateUserCommand) (user.Model, error)
}

func (s UserService) CreateUser(command user.CreateUserCommand, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserCreatePermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	// Validate user
	// Map command to user
	userModel, err := user.NewUser(primitive.NewObjectID().Hex(), command.Username,
		command.Password, command.PersonInfo, command.ContactInfo, time.Now(), 1)
	if err := userModel.Validate(); err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}
	userExists := s.userRepository.ExistsByUsernameAndEmail(userModel.Username, userModel.ContactInfo.Email)

	if userExists {
		return user.Model{}, constants.ErrorConflict
	}

	// Hash the password
	err = hashPassword(&userModel)
	if err != nil {
		return user.Model{}, err
	}

	savedUser, err := s.userRepository.SaveUser(userModel)
	if err != nil {
		return user.Model{}, err
	}

	// Add default roles to user
	err = s.addDefaultRoles(savedUser.ID, authContext)

	userAfterCreate, err := s.userRepository.GetUser(savedUser.ID)
	if err != nil {
		return user.Model{}, err
	}
	log.Log("User %s created successfully", userAfterCreate.Username)
	return userAfterCreate, nil
}

// addDefaultRoles adds default roles to a user
func (s UserService) addDefaultRoles(userId primitive.ObjectID, authContext auth.PermissionContext) error {
	// Get all default roles from all registered role providers
	defaultRoleNames := manager.GetAllDefaultRoleNames()
	if len(defaultRoleNames) == 0 {
		log.Log("No default roles found")
		return nil
	}
	for _, roleName := range defaultRoleNames {
		roleModel, err := s.roleService.GetRoleByName(roleName, authContext)
		if err != nil {
			return err
		}
		// Add roles to user
		err = s.userRepository.AddRoleToUser(userId, roleModel.ID)
	}

	return nil
}

// GetUsers gets all users
func (s UserService) GetUsers(authContext auth.PermissionContext) ([]user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return nil, constants.ErrorPermissionDenied
	}
	return s.userRepository.GetUsers()

}

func (s UserService) GetUser(id string, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}
	userModel, err := s.userRepository.GetUser(objID)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

func (s UserService) ExistsUser(id string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return false, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, constants.ErrorBadRequest
	}
	exists, err := s.userRepository.ExistsUser(objID)
	if err != nil {
		return false, constants.ErrorInternalServerError
	}
	return exists, nil
}

func (s UserService) DeleteUser(command user.DeleteUserCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserDeletePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	existsUser, err := s.ExistsUser(command.ID.Hex(), authContext)
	if err != nil {
		return constants.ErrorInternalServerError
	}
	if !existsUser {
		return constants.ErrorNotFound
	}

	err = s.userRepository.DeleteUser(command.ID)
	if err != nil {
		return constants.ErrorInternalServerError
	}
	return nil
}

// hashPassword hashes a password using bcrypt and assigns it to the user
func hashPassword(userModel *user.Model) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), bcrypt.DefaultCost)
	if err != nil {
		return constants.ErrorInternalServerError
	}
	userModel.Password = string(hashedPassword)
	return nil
}

// VerifyUser verifies a user by username and password
func (s UserService) VerifyUser(username string, password string, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	// Get the user by username
	userModel, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	// Compare the passwords
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(password))
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return userModel, nil
}

// ExistsByUsername gets a user by username
func (s UserService) ExistsByUsername(username string) (bool, error) {
	return s.userRepository.ExistsByUsername(username), nil
}

// AddRoleToUser adds a role to a user
func (s UserService) AddRoleToUser(command user.AddRoleToUserCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserUpdatePermission) != nil {
		return constants.ErrorPermissionDenied
	}
	return s.userRepository.AddRoleToUser(command.UserID, command.RoleID)
}

// RemoveRoleFromUser removes a role from a user
func (s UserService) RemoveRoleFromUser(command user.RemoveRoleFromUserCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserDeletePermission) != nil {
		return constants.ErrorPermissionDenied
	}
	return s.userRepository.RemoveRoleFromUser(command.UserID, command.RoleID)
}

// GetUserRoles gets the roles of a user
func (s UserService) GetUserRoles(userID primitive.ObjectID, authContext auth.PermissionContext) ([]role.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return nil, constants.ErrorPermissionDenied
	}
	return s.userRepository.GetUserRoles(userID)
}

// GetUserPermissionList gets the permissionList of a user
func (s UserService) GetUserPermissionList(userID primitive.ObjectID) ([]auth.Permission, error) {
	userRoles, err := s.GetUserRoles(userID, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})
	if err != nil {
		return nil, constants.ErrorInternalServerError
	}

	var userPermissionList []auth.Permission
	for _, userRole := range userRoles {
		userPermissionList = append(userPermissionList, userRole.Permissions...)
	}

	return userPermissionList, nil
}

// UpdateUser updates a user
func (s UserService) UpdateUser(id string, command user.UpdateUserCommand, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserUpdatePermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}

	// Update user
	userModel, err := s.userRepository.UpdateUser(objID, command)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

// UpdateUserSelf updates a user by itself
func (s UserService) UpdateUserSelf(command user.UpdateUserCommand, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserSelfUpdatePermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(authContext.UserId.Hex())
	if err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}

	// Update user
	userModel, err := s.userRepository.UpdateUser(objID, command)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}
