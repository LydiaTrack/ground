package service

import (
	"github.com/LydiaTrack/ground/pkg/responses"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/LydiaTrack/ground/internal/log"
	"github.com/LydiaTrack/ground/internal/permissions"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/registry"
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
	GetUsers() (responses.QueryResult[user.Model], error)
	// GetUser gets a user by id
	GetUser(id primitive.ObjectID) (user.Model, error)
	// GetUserByUsername gets a user by username
	GetUserByUsername(username string) (user.Model, error)
	// ExistsUser checks if a user exists
	ExistsUser(id primitive.ObjectID) (bool, error)
	// DeleteUser deletes a user by id
	DeleteUser(id primitive.ObjectID) error
	// ExistsByUsernameAndEmail checks if a user exists by username
	ExistsByUsernameAndEmail(username string, email string) bool
	// ExistsByUsername checks if a user exists by username
	ExistsByUsername(username string) bool
	// ExistsByEmail checks if a user exists by email address
	ExistsByEmail(email string) bool
	// AddRoleToUser adds a role to a user
	AddRoleToUser(userID primitive.ObjectID, roleID primitive.ObjectID) error
	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(userID primitive.ObjectID, roleID primitive.ObjectID) error
	// GetUserRoles gets the roles of a user
	GetUserRoles(userID primitive.ObjectID) (responses.QueryResult[role.Model], error)
	// UpdateUser updates a user and returns the updated user
	UpdateUser(id primitive.ObjectID, updateCommand user.UpdateUserCommand) (user.Model, error)
	// UpdateUserPassword updates a user's password
	UpdateUserPassword(id primitive.ObjectID, password string) error
	// GetUserByEmailAddress gets a user by email address
	GetUserByEmailAddress(email string) (user.Model, error)
}

func (s UserService) CreateUser(command user.CreateUserCommand, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserCreatePermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	// Validate user
	// Map command to user
	userModel, err := user.NewUser(primitive.NewObjectID().Hex(), command.Username,
		command.Password, command.PersonInfo, command.ContactInfo, time.Now(), command.Properties, 1)
	if err := userModel.Validate(); err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}
	userExists := s.userRepository.ExistsByUsernameAndEmail(userModel.Username, userModel.ContactInfo.Email)

	if userExists {
		return user.Model{}, constants.ErrorConflict
	}

	// Hash the password
	err = hashUserPassword(&userModel)
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
	defaultRoleNames := registry.GetAllDefaultRoleNames()
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
func (s UserService) GetUsers(authContext auth.PermissionContext) (responses.QueryResult[user.Model], error) {
	// Check if the user has the required permission
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		// Return a pointer to an empty QueryResult and the permission denied error
		return responses.QueryResult[user.Model]{}, constants.ErrorPermissionDenied
	}

	// Fetch users from the repository
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

func (s UserService) GetSelfUser(authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserSelfGetPermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(authContext.UserId.Hex())
	if err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}
	userModel, err := s.userRepository.GetUser(objID)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

func (s UserService) GetUserByEmail(email string, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	userModel, err := s.userRepository.GetUserByEmailAddress(email)
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

// hashUserPassword hashes a password using bcrypt and assigns it to the user
func hashUserPassword(userModel *user.Model) error {
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

// ExistsByUsername checks whether a user exists with specified username
func (s UserService) ExistsByUsername(username string) bool {
	return s.userRepository.ExistsByUsername(username)
}

// ExistsByEmail checks whether a user exists with specified email address
func (s UserService) ExistsByEmail(email string) bool {
	return s.userRepository.ExistsByEmail(email)
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
func (s UserService) GetUserRoles(userID primitive.ObjectID, authContext auth.PermissionContext) (responses.QueryResult[role.Model], error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return responses.QueryResult[role.Model]{}, constants.ErrorPermissionDenied
	}
	return s.userRepository.GetUserRoles(userID)
}

// GetUserPermissionList gets the permissionList of a user
func (s UserService) GetUserPermissionList(userID primitive.ObjectID) ([]auth.Permission, error) {
	result, err := s.GetUserRoles(userID, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      nil,
	})
	if err != nil {
		return nil, constants.ErrorInternalServerError
	}

	var userPermissionList []auth.Permission
	for _, userRole := range result.Data {
		userPermissionList = append(userPermissionList, userRole.Permissions...)
	}

	return userPermissionList, nil
}

// UpdateUser updates a user
func (s UserService) UpdateUser(id string, command user.UpdateUserCommand, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserUpdatePermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	// Validate UpdateUserCommand
	if err := command.Validate(); err != nil {
		return user.Model{}, constants.ErrorBadRequest
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

	// Validate UpdateUserCommand
	if err := command.Validate(); err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}

	objID, err := primitive.ObjectIDFromHex(authContext.UserId.Hex())
	if err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}

	userToUpdate, err := s.GetSelfUser(authContext)
	if err != nil {
		return user.Model{}, err
	}

	// Check if the user is trying to update another user's email
	if userToUpdate.ContactInfo.Email != command.ContactInfo.Email {
		// If the email is different, check if the new email is already in use
		if s.ExistsByEmail(command.ContactInfo.Email) {
			return user.Model{}, constants.ErrorConflict
		}
	}

	// Update user
	userModel, err := s.userRepository.UpdateUser(objID, command)
	if err != nil {
		return user.Model{}, err
	}
	return userModel, nil
}

// UpdateUserPassword updates a user's password
func (s UserService) UpdateUserPassword(id string, cmd user.UpdatePasswordCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserSelfUpdatePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	// Check if the user is trying to update another user's password
	if authContext.UserId.Hex() != id {
		return constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.ErrorBadRequest
	}

	userModel, err := s.GetUser(id, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      &objID,
	})
	if err != nil {
		return err
	}

	// Check if the current password is correct
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(cmd.CurrentPassword))
	if err != nil {
		return constants.ErrorUnauthorized
	}

	// Assign the new password, theres no need to store old passwords
	userModel.Password = cmd.NewPassword

	// Hash the password
	err = hashUserPassword(&userModel)
	if err != nil {
		return err
	}

	// Update password
	err = s.userRepository.UpdateUserPassword(objID, userModel.Password)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserPasswordSelf updates a user's password by itself
func (s UserService) UpdateUserPasswordSelf(cmd user.UpdatePasswordCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserSelfUpdatePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	objID, err := primitive.ObjectIDFromHex(authContext.UserId.Hex())
	if err != nil {
		return constants.ErrorBadRequest
	}

	userModel, err := s.GetUser(authContext.UserId.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      &objID,
	})
	if err != nil {
		return err
	}

	// Check if the current password is correct
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(cmd.CurrentPassword))
	if err != nil {
		return constants.ErrorBadRequest
	}

	// Assign the new password, theres no need to store old passwords
	userModel.Password = cmd.NewPassword

	// Hash the password
	err = hashUserPassword(&userModel)
	if err != nil {
		return err
	}

	// Update password
	err = s.userRepository.UpdateUserPassword(objID, userModel.Password)
	if err != nil {
		return err
	}
	return nil

}

// ResetUserPassword resets a user's password without knowing the current password
func (s UserService) ResetUserPassword(id string, cmd user.ResetPasswordCommand) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.ErrorBadRequest
	}

	userModel, err := s.GetUser(id, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserId:      &objID,
	})
	if err != nil {
		return err
	}

	userModel.Password = cmd.NewPassword

	// Hash the password
	err = hashUserPassword(&userModel)
	if err != nil {
		return err
	}

	// Update password
	err = s.userRepository.UpdateUserPassword(objID, userModel.Password)
	if err != nil {
		return err
	}
	return nil
}
