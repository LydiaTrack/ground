package service

import (
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatsRepository interface {
	GetStatsByUserID(userID primitive.ObjectID) (user.StatsDocument, error)
	CreateStats(stats user.StatsDocument) error
	UpdateStats(stats user.StatsDocument) error
	IncrementField(statsID primitive.ObjectID, fieldName string, increment int) error
	IncrementInt64Field(statsID primitive.ObjectID, fieldName string, increment int64) error
	UpdateField(statsID primitive.ObjectID, fieldName string, value interface{}) error
	UpdateFields(statsID primitive.ObjectID, fields map[string]interface{}) error
}

type UserStatsService struct {
	userStatsRepository UserStatsRepository
}

func NewUserStatsService(userStatsRepository UserStatsRepository) *UserStatsService {
	return &UserStatsService{
		userStatsRepository: userStatsRepository,
	}
}

// CreateUserStats creates a new stats document for a user
func (s *UserStatsService) CreateUserStats(userID primitive.ObjectID, username string) error {
	stats := user.NewStatsDocument(userID, username)
	if err := s.userStatsRepository.CreateStats(stats); err != nil {
		log.Log("Failed to create user stats: %v", err)
		return constants.ErrorInternalServerError
	}
	return nil
}

// GetUserStats retrieves a user's stats
func (s *UserStatsService) GetUserStats(userID primitive.ObjectID, authContext auth.PermissionContext) (user.StatsDocument, error) {
	// Check if the requesting user is the same as the user ID in the stats or has admin permission
	if authContext.UserID != nil && *authContext.UserID != userID && !auth.HasPermission(authContext.Permissions, auth.AdminPermission) {
		return user.StatsDocument{}, constants.ErrorPermissionDenied
	}

	stats, err := s.userStatsRepository.GetStatsByUserID(userID)
	if err != nil {
		return user.StatsDocument{}, constants.ErrorInternalServerError
	}
	return stats, nil
}

// UpdateUserStats updates a user's stats
func (s *UserStatsService) UpdateUserStats(stats user.StatsDocument, authContext auth.PermissionContext) error {
	// Check if the requesting user is the same as the user ID in the stats or has admin permission
	core := stats.GetCoreFields()
	if authContext.UserID != nil && *authContext.UserID != core.UserID && !auth.HasPermission(authContext.Permissions, auth.AdminPermission) {
		return constants.ErrorPermissionDenied
	}

	// Calculate general stat fields before updating
	stats.CalculateStatFields()

	if err := s.userStatsRepository.UpdateStats(stats); err != nil {
		return constants.ErrorInternalServerError
	}
	return nil
}

// RecordLogin increments login count and updates last login date
func (s *UserStatsService) RecordLogin(userID primitive.ObjectID, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	fields := make(map[string]interface{})
	fields["totalLogins"] = stats.GetInt("totalLogins") + 1
	core := stats.GetCoreFields()

	return s.userStatsRepository.UpdateFields(core.ID, fields)
}

// IncrementField increments a numeric field in the user's stats
func (s *UserStatsService) IncrementField(userID primitive.ObjectID, fieldName string, increment int, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	core := stats.GetCoreFields()
	return s.userStatsRepository.IncrementField(core.ID, fieldName, increment)
}

// IncrementInt64Field increments a numeric int64 field in the user's stats
func (s *UserStatsService) IncrementInt64Field(userID primitive.ObjectID, fieldName string, increment int64, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	core := stats.GetCoreFields()
	return s.userStatsRepository.IncrementInt64Field(core.ID, fieldName, increment)
}

// UpdateField updates a specific field in the user's stats
func (s *UserStatsService) UpdateField(userID primitive.ObjectID, fieldName string, value interface{}, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	core := stats.GetCoreFields()
	return s.userStatsRepository.UpdateField(core.ID, fieldName, value)
}

// UpdateFields updates multiple fields in the user's stats
func (s *UserStatsService) UpdateFields(userID primitive.ObjectID, fields map[string]interface{}, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	core := stats.GetCoreFields()
	return s.userStatsRepository.UpdateFields(core.ID, fields)
}
