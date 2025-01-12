package service

import (
	"context"
	"github.com/LydiaTrack/ground/internal/log"
	"github.com/LydiaTrack/ground/internal/permissions"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/role"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/mongodb/repository"
	"github.com/LydiaTrack/ground/pkg/registry"
	"github.com/LydiaTrack/ground/pkg/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var userSearchFields = []string{"username", "contactInfo.email"}

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
	repository.Repository[user.Model]
	ExistsByUsernameAndEmail(username, email string) bool
	ExistsByUsername(username string) bool
	ExistsByEmail(email string) bool
	GetByUsername(username string) (user.Model, error)
	GetByEmail(email string) (user.Model, error)
	AddRole(userID, roleID primitive.ObjectID) error
	RemoveRole(userID, roleID primitive.ObjectID) error
	GetUserRoles(userID primitive.ObjectID) (responses.QueryResult[role.Model], error)
	UpdateUserPassword(id primitive.ObjectID, password string) error
}

// Create creates a new user
func (s UserService) Create(command user.CreateUserCommand, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserCreatePermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	userModel, err := user.NewUser(
		user.WithUsername(command.Username),
		user.WithPassword(command.Password),
		user.WithContactInfo(command.ContactInfo),
		user.WithPersonInfo(command.PersonInfo),
		user.WithProperties(command.Properties),
	)

	if err := userModel.Validate(); err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}

	if s.userRepository.ExistsByUsernameAndEmail(userModel.Username, userModel.ContactInfo.Email) {
		return user.Model{}, constants.ErrorConflict
	}

	if err = hashUserPassword(userModel); err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	createResult, err := s.userRepository.Create(context.Background(), *userModel)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	insertedId := createResult.InsertedID.(primitive.ObjectID)

	err = s.addDefaultRoles(insertedId, authContext)
	if err != nil {
		return user.Model{}, err
	}

	savedUser, err := s.userRepository.GetByID(context.Background(), insertedId)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return savedUser, nil
}

// InitializeDefaultRolesForAllUsers initializes default roles for all users
func (s UserService) InitializeDefaultRolesForAllUsers() error {
	// Retrieve all users from the repository which does not have any default roles
	defaultRoleIDs, err := s.getDefaultRoleIDs()
	if err != nil {
		return err
	}

	if defaultRoleIDs == nil {
		log.Log("No default roles found in the system. Skipping default role assignment.")
		return nil
	}

	filter := bson.M{"roleIds": bson.M{"$nin": defaultRoleIDs}}
	allUsers, err := s.userRepository.Query(context.Background(), filter, nil, "")
	if err != nil {
		return err
	}

	// Create a permission context that grants all necessary permissions.
	// Here, we assume the AdminPermission allows adding roles to users.
	adminCtx := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	if allUsers.TotalElements == 0 {
		log.Log("No users found in the system without default roles. Skipping default role assignment.")
		return nil
	}

	// Iterate over each user and attempt to add default roles
	for _, usr := range allUsers.Data {
		err := s.addDefaultRoles(usr.ID, adminCtx)
		if err != nil {
			// Log the error but continue with other users
			log.Log("Failed to add default roles to user %s: %v", usr.Username, err)
		} else {
			log.Log("Default roles successfully assigned to user %s", usr.Username)
		}
	}

	return nil
}

// Query users by an optional search text
func (s UserService) Query(searchText string, authContext auth.PermissionContext) (responses.QueryResult[user.Model], error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return responses.QueryResult[user.Model]{}, constants.ErrorPermissionDenied
	}

	return s.userRepository.Query(context.Background(), nil, userSearchFields, searchText)
}

// QueryPaginated query users by an optional search text with pagination
func (s UserService) QueryPaginated(searchText string, page, limit int, authContext auth.PermissionContext) (repository.PaginatedResult[user.Model], error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return repository.PaginatedResult[user.Model]{}, constants.ErrorPermissionDenied
	}

	ctx := context.Background()
	return s.userRepository.QueryPaginate(ctx, nil, userSearchFields, searchText, page, limit, nil)
}

// GetByUsername retrieves a user by username
func (s UserService) GetByUsername(username string, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	userModel, err := s.userRepository.GetByUsername(username)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return userModel, nil
}

// GetByEmail retrieves a user by email
func (s UserService) GetByEmail(email string, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	userModel, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return userModel, nil
}

// Get retrieves a user by ID
func (s UserService) Get(id string, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	userModel, err := s.userRepository.GetByID(context.Background(), id)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return userModel, nil
}

// Exists checks if a user exists by ID
func (s UserService) Exists(id string) (bool, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, constants.ErrorBadRequest
	}

	exists, err := s.userRepository.ExistsByID(context.Background(), objID)
	if err != nil {
		return false, constants.ErrorInternalServerError
	}

	return exists, nil
}

// ExistsByUsername checks if a user exists by username
func (s UserService) ExistsByUsername(username string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return false, constants.ErrorPermissionDenied
	}

	exists := s.userRepository.ExistsByUsername(username)
	return exists, nil
}

// ExistsByEmail checks if a user exists by email
func (s UserService) ExistsByEmail(email string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return false, constants.ErrorPermissionDenied
	}

	exists := s.userRepository.ExistsByEmail(email)
	return exists, nil
}

// ExistsByEmailAndUsername checks if a user exists by email and username
func (s UserService) ExistsByEmailAndUsername(email string, username string, authContext auth.PermissionContext) (bool, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return false, constants.ErrorPermissionDenied
	}

	existsByEmail := s.userRepository.ExistsByEmail(email)
	existsByUsername := s.userRepository.ExistsByUsername(username)
	return existsByEmail || existsByUsername, nil
}

// Delete deletes a user by ID
func (s UserService) Delete(command user.DeleteUserCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserDeletePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	exists, err := s.Exists(command.ID.Hex())
	if err != nil {
		return constants.ErrorNotFound
	}

	if !exists {
		return constants.ErrorNotFound
	}

	_, err = s.userRepository.Delete(context.Background(), command.ID)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	return nil
}

// Update updates a user
func (s UserService) Update(id string, command user.UpdateUserCommand, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserUpdatePermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	exists, err := s.Exists(id)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}
	if !exists {
		return user.Model{}, constants.ErrorNotFound
	}

	if err := command.Validate(); err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}

	_, err = s.userRepository.Update(context.Background(), objID, command)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	updatedUser, err := s.userRepository.GetByID(context.Background(), objID)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return updatedUser, nil
}

// UpdateSelf updates a user's own information
func (s UserService) UpdateSelf(command user.UpdateUserCommand, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserSelfUpdatePermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	exists, err := s.Exists(authContext.UserID.Hex())
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}
	if !exists {
		return user.Model{}, constants.ErrorNotFound
	}

	if err := command.Validate(); err != nil {
		return user.Model{}, constants.ErrorBadRequest
	}

	_, err = s.userRepository.Update(context.Background(), *authContext.UserID, command)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	updatedUser, err := s.userRepository.GetByID(context.Background(), *authContext.UserID)
	if err != nil {
		return user.Model{}, constants.ErrorInternalServerError
	}

	return updatedUser, nil
}

// UpdatePassword updates a user's password
func (s UserService) UpdatePassword(id string, cmd user.UpdatePasswordCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserSelfUpdatePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	userModel, err := s.Get(id, authContext)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	// Verify current password
	if err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(cmd.CurrentPassword)); err != nil {
		return constants.ErrorUnauthorized
	}

	// Hash and update the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.ErrorBadRequest
	}

	err = s.userRepository.UpdateUserPassword(objID, string(hashedPassword))
	if err != nil {
		return constants.ErrorInternalServerError
	}

	return nil
}

// UpdateSelfPassword updates a user's own password
func (s UserService) UpdateSelfPassword(command user.UpdatePasswordCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserSelfUpdatePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	userModel, err := s.Get(authContext.UserID.Hex(), authContext)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	// Verify current password
	if err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(command.CurrentPassword)); err != nil {
		return constants.ErrorUnauthorized
	}

	// Hash and update the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(command.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	err = s.userRepository.UpdateUserPassword(*authContext.UserID, string(hashedPassword))
	if err != nil {
		return constants.ErrorInternalServerError
	}

	return nil
}

// ResetPassword resets a user's password without knowing the current password
func (s UserService) ResetPassword(id string, cmd user.ResetPasswordCommand) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.ErrorBadRequest
	}

	userModel, err := s.Get(id, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      &objID,
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

// VerifyUser verifies a user by username and password
func (s UserService) VerifyUser(username, password string, authContext auth.PermissionContext) (user.Model, error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return user.Model{}, constants.ErrorPermissionDenied
	}

	userModel, err := s.userRepository.GetByUsername(username)
	if err != nil {
		return user.Model{}, constants.ErrorNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(password)); err != nil {
		return user.Model{}, constants.ErrorUnauthorized
	}

	return userModel, nil
}

// GetRoles retrieves roles for a user
func (s UserService) GetRoles(userID primitive.ObjectID, authContext auth.PermissionContext) (responses.QueryResult[role.Model], error) {
	if auth.CheckPermission(authContext.Permissions, permissions.UserReadPermission) != nil {
		return responses.QueryResult[role.Model]{}, constants.ErrorPermissionDenied
	}

	return s.userRepository.GetUserRoles(userID)
}

// AddRole adds a role to a user
func (s UserService) AddRole(command user.AddRoleToUserCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserUpdatePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	exists, err := s.roleService.Exists(command.RoleID.Hex(), authContext)
	if err != nil {
		return constants.ErrorInternalServerError
	}
	if !exists {
		return constants.ErrorNotFound
	}

	err = s.userRepository.AddRole(command.UserID, command.RoleID)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	return nil
}

// RemoveRole removes a role from a user
func (s UserService) RemoveRole(command user.RemoveRoleFromUserCommand, authContext auth.PermissionContext) error {
	if auth.CheckPermission(authContext.Permissions, permissions.UserUpdatePermission) != nil {
		return constants.ErrorPermissionDenied
	}

	exists, err := s.roleService.Exists(command.RoleID.Hex(), authContext)
	if err != nil {
		return constants.ErrorInternalServerError
	}
	if !exists {
		return constants.ErrorNotFound
	}

	err = s.userRepository.RemoveRole(command.UserID, command.RoleID)
	if err != nil {
		return constants.ErrorInternalServerError
	}

	return nil
}

// GetPermissionList retrieves permissions for a user
func (s UserService) GetPermissionList(userID primitive.ObjectID) ([]auth.Permission, error) {
	userRoles, err := s.GetRoles(userID, auth.CreateAdminAuthContext())
	if err != nil {
		return nil, err
	}

	var permissionList []auth.Permission
	for _, roleModel := range userRoles.Data {
		permissionList = append(permissionList, roleModel.Permissions...)
	}

	return permissionList, nil
}

// addDefaultRoles adds default roles to a user
func (s UserService) addDefaultRoles(userID primitive.ObjectID, authContext auth.PermissionContext) error {
	defaultRoleIDs, err := s.getDefaultRoleIDs()
	if err != nil {
		return err
	}

	for _, id := range defaultRoleIDs {

		if err = s.userRepository.AddRole(userID, id); err != nil {
			return err
		}
	}

	return nil
}

// getDefaultRoleIDs retrieves default role IDs
func (s UserService) getDefaultRoleIDs() ([]primitive.ObjectID, error) {
	defaultRoleNames := registry.GetAllDefaultRoleNames()
	if len(defaultRoleNames) == 0 {
		return nil, nil
	}

	var roleIDs []primitive.ObjectID
	for _, roleName := range defaultRoleNames {
		roleModel, err := s.roleService.GetByName(roleName, auth.CreateAdminAuthContext())
		if err != nil {
			return nil, err
		}

		roleIDs = append(roleIDs, roleModel.ID)
	}

	return roleIDs, nil
}

// hashUserPassword hashes the user's password
func hashUserPassword(userModel *user.Model) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userModel.Password = string(hashedPassword)
	return nil
}
