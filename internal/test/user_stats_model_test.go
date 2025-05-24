package test

import (
	"testing"
	"time"

	"github.com/LydiaTrack/ground/pkg/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserStatsModel(t *testing.T) {
	t.Run("NewStats", testNewStats)
	t.Run("CalculateStatFields", testCalculateStatFields)
	t.Run("CalculateStatFieldsNewDay", testCalculateStatFieldsNewDay)
	t.Run("CalculateStatFieldsSameDay", testCalculateStatFieldsSameDay)
}

func testNewStats(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Record the time before creating stats
	beforeCreate := time.Now()

	stats := user.NewStats(userID, username)

	// Record the time after creating stats
	afterCreate := time.Now()

	// Verify basic fields
	if stats.UserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, stats.UserID)
	}

	if stats.Username != username {
		t.Errorf("Expected username %s, got %s", username, stats.Username)
	}

	// Verify timestamps are reasonable
	if stats.CreatedDate.Before(beforeCreate) || stats.CreatedDate.After(afterCreate) {
		t.Errorf("CreatedDate %v should be between %v and %v", stats.CreatedDate, beforeCreate, afterCreate)
	}

	if stats.UpdatedDate.Before(beforeCreate) || stats.UpdatedDate.After(afterCreate) {
		t.Errorf("UpdatedDate %v should be between %v and %v", stats.UpdatedDate, beforeCreate, afterCreate)
	}

	if stats.LastActiveDate.Before(beforeCreate) || stats.LastActiveDate.After(afterCreate) {
		t.Errorf("LastActiveDate %v should be between %v and %v", stats.LastActiveDate, beforeCreate, afterCreate)
	}

	// Verify initial values
	if stats.TotalLogins != 1 {
		t.Errorf("Expected initial TotalLogins to be 1, got %d", stats.TotalLogins)
	}

	if stats.ActiveDaysCount != 1 {
		t.Errorf("Expected initial ActiveDaysCount to be 1, got %d", stats.ActiveDaysCount)
	}

	if stats.DayAge != 0 {
		t.Errorf("Expected initial DayAge to be 0, got %d", stats.DayAge)
	}

	// Verify all other stats start at 0
	if stats.TasksCreated != 0 {
		t.Errorf("Expected initial TasksCreated to be 0, got %d", stats.TasksCreated)
	}

	if stats.TasksCompleted != 0 {
		t.Errorf("Expected initial TasksCompleted to be 0, got %d", stats.TasksCompleted)
	}

	if stats.NotesCreated != 0 {
		t.Errorf("Expected initial NotesCreated to be 0, got %d", stats.NotesCreated)
	}

	if stats.TotalTimeTracked != 0 {
		t.Errorf("Expected initial TotalTimeTracked to be 0, got %d", stats.TotalTimeTracked)
	}

	if stats.TimeEntryCount != 0 {
		t.Errorf("Expected initial TimeEntryCount to be 0, got %d", stats.TimeEntryCount)
	}

	if stats.ProjectsCreated != 0 {
		t.Errorf("Expected initial ProjectsCreated to be 0, got %d", stats.ProjectsCreated)
	}
}

func testCalculateStatFields(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with a past creation date
	pastTime := time.Now().Add(-48 * time.Hour) // 2 days ago
	stats := &user.StatsModel{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Username:        username,
		CreatedDate:     pastTime,
		UpdatedDate:     pastTime,
		LastActiveDate:  pastTime,
		ActiveDaysCount: 1,
		DayAge:          0,
	}

	// Record time before calculation
	beforeCalculation := time.Now()

	// Call CalculateStatFields
	stats.CalculateStatFields()

	// Record time after calculation
	afterCalculation := time.Now()

	// Verify UpdatedDate was updated
	if stats.UpdatedDate.Before(beforeCalculation) || stats.UpdatedDate.After(afterCalculation) {
		t.Errorf("UpdatedDate should be updated to current time, got %v", stats.UpdatedDate)
	}

	// Verify LastActiveDate was updated
	if stats.LastActiveDate.Before(beforeCalculation) || stats.LastActiveDate.After(afterCalculation) {
		t.Errorf("LastActiveDate should be updated to current time, got %v", stats.LastActiveDate)
	}

	// Verify DayAge was calculated (should be around 2 days)
	expectedDayAge := int(time.Now().Sub(stats.CreatedDate).Hours() / 24)
	if stats.DayAge < expectedDayAge-1 || stats.DayAge > expectedDayAge+1 {
		t.Errorf("Expected DayAge to be around %d, got %d", expectedDayAge, stats.DayAge)
	}

	// ActiveDaysCount should be incremented since it's a new day
	if stats.ActiveDaysCount != 2 {
		t.Errorf("Expected ActiveDaysCount to be incremented to 2, got %d", stats.ActiveDaysCount)
	}
}

func testCalculateStatFieldsNewDay(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with yesterday's last active date
	yesterday := time.Now().Add(-24 * time.Hour)
	stats := &user.StatsModel{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Username:        username,
		CreatedDate:     yesterday,
		UpdatedDate:     yesterday,
		LastActiveDate:  yesterday,
		ActiveDaysCount: 1,
		DayAge:          0,
	}

	initialActiveDaysCount := stats.ActiveDaysCount

	// Call CalculateStatFields
	stats.CalculateStatFields()

	// ActiveDaysCount should be incremented for new day
	if stats.ActiveDaysCount != initialActiveDaysCount+1 {
		t.Errorf("Expected ActiveDaysCount to be incremented from %d to %d, got %d",
			initialActiveDaysCount, initialActiveDaysCount+1, stats.ActiveDaysCount)
	}
}

func testCalculateStatFieldsSameDay(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with today's date (a few hours ago)
	todayEarlier := time.Now().Add(-2 * time.Hour)
	stats := &user.StatsModel{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Username:        username,
		CreatedDate:     todayEarlier,
		UpdatedDate:     todayEarlier,
		LastActiveDate:  todayEarlier,
		ActiveDaysCount: 1,
		DayAge:          0,
	}

	initialActiveDaysCount := stats.ActiveDaysCount

	// Call CalculateStatFields
	stats.CalculateStatFields()

	// ActiveDaysCount should NOT be incremented for same day
	if stats.ActiveDaysCount != initialActiveDaysCount {
		t.Errorf("Expected ActiveDaysCount to remain %d for same day activity, got %d",
			initialActiveDaysCount, stats.ActiveDaysCount)
	}

	// But other fields should still be updated
	now := time.Now()
	if stats.LastActiveDate.Before(now.Add(-1*time.Minute)) || stats.LastActiveDate.After(now.Add(1*time.Minute)) {
		t.Errorf("LastActiveDate should be updated to current time")
	}

	if stats.UpdatedDate.Before(now.Add(-1*time.Minute)) || stats.UpdatedDate.After(now.Add(1*time.Minute)) {
		t.Errorf("UpdatedDate should be updated to current time")
	}
}

func TestStatsModelEdgeCases(t *testing.T) {
	t.Run("ZeroLastActiveDate", testZeroLastActiveDate)
	t.Run("FutureCreatedDate", testFutureCreatedDate)
}

func testZeroLastActiveDate(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with zero LastActiveDate
	stats := &user.StatsModel{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Username:        username,
		CreatedDate:     time.Now().Add(-24 * time.Hour),
		UpdatedDate:     time.Now().Add(-24 * time.Hour),
		LastActiveDate:  time.Time{}, // Zero value
		ActiveDaysCount: 1,
		DayAge:          0,
	}

	initialActiveDaysCount := stats.ActiveDaysCount

	// Call CalculateStatFields
	stats.CalculateStatFields()

	// ActiveDaysCount should NOT be incremented when LastActiveDate is zero
	if stats.ActiveDaysCount != initialActiveDaysCount {
		t.Errorf("Expected ActiveDaysCount to remain %d when LastActiveDate is zero, got %d",
			initialActiveDaysCount, stats.ActiveDaysCount)
	}

	// LastActiveDate should be set to current time
	now := time.Now()
	if stats.LastActiveDate.Before(now.Add(-1*time.Minute)) || stats.LastActiveDate.After(now.Add(1*time.Minute)) {
		t.Errorf("LastActiveDate should be set to current time when it was zero")
	}
}

func testFutureCreatedDate(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with future creation date (edge case)
	futureTime := time.Now().Add(24 * time.Hour)
	stats := &user.StatsModel{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Username:        username,
		CreatedDate:     futureTime,
		UpdatedDate:     futureTime,
		LastActiveDate:  futureTime,
		ActiveDaysCount: 1,
		DayAge:          0,
	}

	// Call CalculateStatFields
	stats.CalculateStatFields()

	// DayAge calculation uses integer division, so it will be 0 for differences less than 24 hours
	// For a future date that's exactly 24 hours away, the calculation would be:
	// int(-24 hours / 24) = int(-1) = -1
	// But due to precision and timing, it might be 0
	// Let's test that it's <= 0 instead of strictly negative
	if stats.DayAge > 0 {
		t.Errorf("Expected DayAge to be <= 0 for future creation date, got %d", stats.DayAge)
	}
}
