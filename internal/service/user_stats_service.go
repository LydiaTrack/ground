package service

import (
	"time"

	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatsRepository interface {
	GetStatsByUserID(userID primitive.ObjectID) (user.StatsModel, error)
	CreateStats(stats *user.StatsModel) error
	UpdateStats(stats *user.StatsModel) error
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
	stats := user.NewStats(userID, username)
	if err := s.userStatsRepository.CreateStats(stats); err != nil {
		log.Log("Failed to create user stats: %v", err)
		return constants.ErrorInternalServerError
	}
	return nil
}

// GetUserStats retrieves a user's stats
func (s *UserStatsService) GetUserStats(userID primitive.ObjectID, authContext auth.PermissionContext) (user.StatsModel, error) {
	// Check if the requesting user is the same as the user ID in the stats or has admin permission
	if authContext.UserID != nil && *authContext.UserID != userID && !auth.HasPermission(authContext.Permissions, auth.AdminPermission) {
		return user.StatsModel{}, constants.ErrorPermissionDenied
	}

	stats, err := s.userStatsRepository.GetStatsByUserID(userID)
	if err != nil {
		return user.StatsModel{}, constants.ErrorInternalServerError
	}
	return stats, nil
}

// UpdateUserStats updates a user's stats
func (s *UserStatsService) UpdateUserStats(stats *user.StatsModel, authContext auth.PermissionContext) error {
	// Check if the requesting user is the same as the user ID in the stats or has admin permission
	if authContext.UserID != nil && *authContext.UserID != stats.UserID && !auth.HasPermission(authContext.Permissions, auth.AdminPermission) {
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
	fields["totalLogins"] = stats.TotalLogins + 1

	return s.userStatsRepository.UpdateFields(stats.ID, fields)
}

// RecordTaskCreated increments the tasks created count
func (s *UserStatsService) RecordTaskCreated(userID primitive.ObjectID, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	return s.userStatsRepository.IncrementField(stats.ID, "tasksCreated", 1)
}

// RecordTaskCompleted increments the tasks completed count
func (s *UserStatsService) RecordTaskCompleted(userID primitive.ObjectID, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	return s.userStatsRepository.IncrementField(stats.ID, "tasksCompleted", 1)
}

// RecordNoteCreated increments the notes created count
func (s *UserStatsService) RecordNoteCreated(userID primitive.ObjectID, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	return s.userStatsRepository.IncrementField(stats.ID, "notesCreated", 1)
}

// RecordTimeEntry adds time tracking data
func (s *UserStatsService) RecordTimeEntry(userID primitive.ObjectID, durationSeconds int64, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	// First increment the time entry count
	if err := s.userStatsRepository.IncrementField(stats.ID, "timeEntryCount", 1); err != nil {
		return err
	}

	// Then update the total time tracked
	return s.userStatsRepository.IncrementInt64Field(stats.ID, "totalTimeTracked", durationSeconds)
}

// RecordProjectCreated increments the projects created count
func (s *UserStatsService) RecordProjectCreated(userID primitive.ObjectID, authContext auth.PermissionContext) error {
	stats, err := s.GetUserStats(userID, authContext)
	if err != nil {
		return err
	}

	return s.userStatsRepository.IncrementField(stats.ID, "projectsCreated", 1)
}

// Helper function to check if two times are on the same day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
