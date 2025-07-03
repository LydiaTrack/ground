package test

import (
	"testing"
	"time"

	"github.com/LydiaTrack/ground/pkg/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserStatsModel(t *testing.T) {
	t.Run("NewStatsDocument", testNewStatsDocument)
	t.Run("CalculateStatFields", testCalculateStatFields)
	t.Run("CalculateStatFieldsNewDay", testCalculateStatFieldsNewDay)
	t.Run("CalculateStatFieldsSameDay", testCalculateStatFieldsSameDay)
}

func testNewStatsDocument(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Record the time before creating stats
	beforeCreate := time.Now()

	stats := user.NewStatsDocument(userID, username)

	// Record the time after creating stats
	afterCreate := time.Now()

	// Verify basic fields
	if stats.GetObjectID("userId") != userID {
		t.Errorf("Expected userID %v, got %v", userID, stats.GetObjectID("userId"))
	}

	if stats.GetString("username") != username {
		t.Errorf("Expected username %s, got %s", username, stats.GetString("username"))
	}

	// Verify timestamps are reasonable
	createdDate := stats.GetTime("createdDate")
	updatedDate := stats.GetTime("updatedDate")
	lastActiveDate := stats.GetTime("lastActiveDate")

	if createdDate.Before(beforeCreate) || createdDate.After(afterCreate) {
		t.Errorf("CreatedDate %v should be between %v and %v", createdDate, beforeCreate, afterCreate)
	}

	if updatedDate.Before(beforeCreate) || updatedDate.After(afterCreate) {
		t.Errorf("UpdatedDate %v should be between %v and %v", updatedDate, beforeCreate, afterCreate)
	}

	if lastActiveDate.Before(beforeCreate) || lastActiveDate.After(afterCreate) {
		t.Errorf("LastActiveDate %v should be between %v and %v", lastActiveDate, beforeCreate, afterCreate)
	}

	// Verify initial values
	if stats.GetInt("totalLogins") != 1 {
		t.Errorf("Expected initial TotalLogins to be 1, got %d", stats.GetInt("totalLogins"))
	}

	if stats.GetInt("activeDaysCount") != 1 {
		t.Errorf("Expected initial ActiveDaysCount to be 1, got %d", stats.GetInt("activeDaysCount"))
	}

	if stats.GetInt("dayAge") != 0 {
		t.Errorf("Expected initial DayAge to be 0, got %d", stats.GetInt("dayAge"))
	}

	// Verify all other stats start at 0 (these fields don't exist initially, so GetInt should return 0)
	if stats.GetInt("tasksCreated") != 0 {
		t.Errorf("Expected initial TasksCreated to be 0, got %d", stats.GetInt("tasksCreated"))
	}

	if stats.GetInt("tasksCompleted") != 0 {
		t.Errorf("Expected initial TasksCompleted to be 0, got %d", stats.GetInt("tasksCompleted"))
	}

	if stats.GetInt("notesCreated") != 0 {
		t.Errorf("Expected initial NotesCreated to be 0, got %d", stats.GetInt("notesCreated"))
	}

	if stats.GetInt64("totalTimeTracked") != 0 {
		t.Errorf("Expected initial TotalTimeTracked to be 0, got %d", stats.GetInt64("totalTimeTracked"))
	}

	if stats.GetInt("timeEntryCount") != 0 {
		t.Errorf("Expected initial TimeEntryCount to be 0, got %d", stats.GetInt("timeEntryCount"))
	}

	if stats.GetInt("projectsCreated") != 0 {
		t.Errorf("Expected initial ProjectsCreated to be 0, got %d", stats.GetInt("projectsCreated"))
	}
}

func testCalculateStatFields(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with a past creation date
	pastTime := time.Now().Add(-48 * time.Hour) // 2 days ago
	stats := user.StatsDocument{
		"_id":             primitive.NewObjectID(),
		"userId":          userID,
		"username":        username,
		"createdDate":     pastTime,
		"updatedDate":     pastTime,
		"lastActiveDate":  pastTime,
		"activeDaysCount": 1,
		"dayAge":          0,
	}

	// Record time before calculation
	beforeCalculation := time.Now()

	// Call CalculateStatFields
	stats.CalculateStatFields()

	// Record time after calculation
	afterCalculation := time.Now()

	// Verify UpdatedDate was updated
	updatedDate := stats.GetTime("updatedDate")
	if updatedDate.Before(beforeCalculation) || updatedDate.After(afterCalculation) {
		t.Errorf("UpdatedDate should be updated to current time, got %v", updatedDate)
	}

	// Verify LastActiveDate was updated
	lastActiveDate := stats.GetTime("lastActiveDate")
	if lastActiveDate.Before(beforeCalculation) || lastActiveDate.After(afterCalculation) {
		t.Errorf("LastActiveDate should be updated to current time, got %v", lastActiveDate)
	}

	// Verify DayAge was calculated (should be around 2 days)
	expectedDayAge := int(time.Now().Sub(stats.GetTime("createdDate")).Hours() / 24)
	dayAge := stats.GetInt("dayAge")
	if dayAge < expectedDayAge-1 || dayAge > expectedDayAge+1 {
		t.Errorf("Expected DayAge to be around %d, got %d", expectedDayAge, dayAge)
	}

	// ActiveDaysCount should be incremented since it's a new day
	if stats.GetInt("activeDaysCount") != 2 {
		t.Errorf("Expected ActiveDaysCount to be incremented to 2, got %d", stats.GetInt("activeDaysCount"))
	}
}

func testCalculateStatFieldsNewDay(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with yesterday's last active date
	yesterday := time.Now().Add(-24 * time.Hour)
	stats := user.StatsDocument{
		"_id":             primitive.NewObjectID(),
		"userId":          userID,
		"username":        username,
		"createdDate":     yesterday,
		"updatedDate":     yesterday,
		"lastActiveDate":  yesterday,
		"activeDaysCount": 1,
		"dayAge":          0,
	}

	initialActiveDaysCount := stats.GetInt("activeDaysCount")

	// Call CalculateStatFields
	stats.CalculateStatFields()

	// ActiveDaysCount should be incremented for new day
	if stats.GetInt("activeDaysCount") != initialActiveDaysCount+1 {
		t.Errorf("Expected ActiveDaysCount to be incremented from %d to %d, got %d",
			initialActiveDaysCount, initialActiveDaysCount+1, stats.GetInt("activeDaysCount"))
	}
}

func testCalculateStatFieldsSameDay(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with today's date (a few hours ago)
	todayEarlier := time.Now().Add(-2 * time.Hour)
	stats := user.StatsDocument{
		"_id":             primitive.NewObjectID(),
		"userId":          userID,
		"username":        username,
		"createdDate":     todayEarlier,
		"updatedDate":     todayEarlier,
		"lastActiveDate":  todayEarlier,
		"activeDaysCount": 1,
		"dayAge":          0,
	}

	initialActiveDaysCount := stats.GetInt("activeDaysCount")

	// Call CalculateStatFields
	stats.CalculateStatFields()

	// ActiveDaysCount should NOT be incremented for same day
	if stats.GetInt("activeDaysCount") != initialActiveDaysCount {
		t.Errorf("Expected ActiveDaysCount to remain %d for same day activity, got %d",
			initialActiveDaysCount, stats.GetInt("activeDaysCount"))
	}

	// But other fields should still be updated
	now := time.Now()
	lastActiveDate := stats.GetTime("lastActiveDate")
	if lastActiveDate.Before(now.Add(-1*time.Minute)) || lastActiveDate.After(now.Add(1*time.Minute)) {
		t.Errorf("LastActiveDate should be updated to current time")
	}
}

func TestStatsModelEdgeCases(t *testing.T) {
	t.Run("ZeroLastActiveDate", testZeroLastActiveDate)
	t.Run("FutureCreatedDate", testFutureCreatedDate)
}

func testZeroLastActiveDate(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with zero lastActiveDate
	now := time.Now()
	stats := user.StatsDocument{
		"_id":             primitive.NewObjectID(),
		"userId":          userID,
		"username":        username,
		"createdDate":     now,
		"updatedDate":     now,
		"lastActiveDate":  time.Time{}, // Zero value
		"activeDaysCount": 1,
		"dayAge":          0,
	}

	// This should not panic or cause issues
	stats.CalculateStatFields()

	// ActiveDaysCount should be incremented since lastActiveDate was zero
	if stats.GetInt("activeDaysCount") != 2 {
		t.Errorf("Expected ActiveDaysCount to be incremented to 2 when lastActiveDate was zero, got %d", stats.GetInt("activeDaysCount"))
	}
}

func testFutureCreatedDate(t *testing.T) {
	userID := primitive.NewObjectID()
	username := "testuser"

	// Create stats with future created date (edge case)
	futureTime := time.Now().Add(24 * time.Hour)
	stats := user.StatsDocument{
		"_id":             primitive.NewObjectID(),
		"userId":          userID,
		"username":        username,
		"createdDate":     futureTime,
		"updatedDate":     futureTime,
		"lastActiveDate":  futureTime,
		"activeDaysCount": 1,
		"dayAge":          0,
	}

	// This should not panic
	stats.CalculateStatFields()

	// DayAge should be negative but handled gracefully
	dayAge := stats.GetInt("dayAge")
	if dayAge > 0 {
		t.Errorf("Expected DayAge to be negative or zero for future created date, got %d", dayAge)
	}
}
