package test

import (
	"errors"
	"testing"

	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/constants"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserStatsRepository is a mock implementation of UserStatsRepository for testing
type MockUserStatsRepository struct {
	stats             map[primitive.ObjectID]user.StatsDocument
	shouldReturnError bool
	errorMessage      string
}

func NewMockUserStatsRepository() *MockUserStatsRepository {
	return &MockUserStatsRepository{
		stats: make(map[primitive.ObjectID]user.StatsDocument),
	}
}

func (m *MockUserStatsRepository) SetError(shouldError bool, message string) {
	m.shouldReturnError = shouldError
	m.errorMessage = message
}

func (m *MockUserStatsRepository) GetStatsByUserID(userID primitive.ObjectID) (user.StatsDocument, error) {
	if m.shouldReturnError {
		return user.StatsDocument{}, errors.New(m.errorMessage)
	}

	for _, stats := range m.stats {
		if stats.GetObjectID("userId") == userID {
			return stats, nil
		}
	}
	return user.StatsDocument{}, errors.New("stats not found")
}

func (m *MockUserStatsRepository) CreateStats(stats user.StatsDocument) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}

	core := stats.GetCoreFields()
	m.stats[core.ID] = stats
	return nil
}

func (m *MockUserStatsRepository) UpdateStats(stats user.StatsDocument) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}

	core := stats.GetCoreFields()
	m.stats[core.ID] = stats
	return nil
}

func (m *MockUserStatsRepository) IncrementField(statsID primitive.ObjectID, fieldName string, increment int) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}

	stats, exists := m.stats[statsID]
	if !exists {
		return errors.New("stats not found")
	}

	// Increment the field using the document method
	stats.IncrementField(fieldName, increment)

	// Update calculated fields
	stats.CalculateStatFields()
	m.stats[statsID] = stats
	return nil
}

func (m *MockUserStatsRepository) IncrementInt64Field(statsID primitive.ObjectID, fieldName string, increment int64) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}

	stats, exists := m.stats[statsID]
	if !exists {
		return errors.New("stats not found")
	}

	// Increment the field using the document method
	stats.IncrementField(fieldName, increment)

	// Update calculated fields
	stats.CalculateStatFields()
	m.stats[statsID] = stats
	return nil
}

func (m *MockUserStatsRepository) UpdateField(statsID primitive.ObjectID, fieldName string, value interface{}) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}

	stats, exists := m.stats[statsID]
	if !exists {
		return errors.New("stats not found")
	}

	// Update the field using the document method
	stats.SetField(fieldName, value)

	// Update calculated fields
	stats.CalculateStatFields()
	m.stats[statsID] = stats
	return nil
}

func (m *MockUserStatsRepository) UpdateFields(statsID primitive.ObjectID, fields map[string]interface{}) error {
	if m.shouldReturnError {
		return errors.New(m.errorMessage)
	}

	stats, exists := m.stats[statsID]
	if !exists {
		return errors.New("stats not found")
	}

	// Update fields using the document method
	for fieldName, value := range fields {
		stats.SetField(fieldName, value)
	}

	// Update calculated fields
	stats.CalculateStatFields()
	m.stats[statsID] = stats
	return nil
}

func TestUserStatsService(t *testing.T) {
	// Initialize logging to prevent nil pointer dereference
	log.InitLogging()

	t.Run("CreateUserStats", testCreateUserStats)
	t.Run("GetUserStats", testGetUserStats)
	t.Run("UpdateUserStats", testUpdateUserStats)
	t.Run("RecordLogin", testRecordLogin)
	t.Run("IncrementField", testIncrementField)
	t.Run("IncrementInt64Field", testIncrementInt64Field)
	t.Run("UpdateField", testUpdateField)
	t.Run("UpdateFields", testUpdateFields)
	t.Run("PermissionChecks", testPermissionChecks)
}

func createTestUserStatsService() (*service.UserStatsService, *MockUserStatsRepository) {
	mockRepo := NewMockUserStatsRepository()
	statsService := service.NewUserStatsService(mockRepo)
	return statsService, mockRepo
}

func testCreateUserStats(t *testing.T) {
	statsService, mockRepo := createTestUserStatsService()

	userID := primitive.NewObjectID()
	username := "testuser"

	// Test successful creation
	err := statsService.CreateUserStats(userID, username)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify stats were created by trying to get the stats
	adminContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	retrievedStats, err := statsService.GetUserStats(userID, adminContext)
	if err != nil {
		t.Errorf("Expected stats to be created and retrievable, got error: %v", err)
	}

	core := retrievedStats.GetCoreFields()
	if core.UserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, core.UserID)
	}

	if core.Username != username {
		t.Errorf("Expected username %s, got %s", username, core.Username)
	}

	// Test creation with repository error
	mockRepo.SetError(true, "repository error")
	err = statsService.CreateUserStats(userID, username)
	if err != constants.ErrorInternalServerError {
		t.Errorf("Expected ErrorInternalServerError, got %v", err)
	}
}

func testGetUserStats(t *testing.T) {
	statsService, mockRepo := createTestUserStatsService()

	userID := primitive.NewObjectID()
	username := "testuser"

	// Create test stats
	stats := user.NewStatsDocument(userID, username)
	mockRepo.CreateStats(stats)

	// Test successful retrieval with admin permission
	adminContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	retrievedStats, err := statsService.GetUserStats(userID, adminContext)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	core := retrievedStats.GetCoreFields()
	if core.UserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, core.UserID)
	}

	// Test retrieval as the same user
	userContext := auth.PermissionContext{
		Permissions: []auth.Permission{},
		UserID:      &userID,
	}

	retrievedStats, err = statsService.GetUserStats(userID, userContext)
	if err != nil {
		t.Errorf("Expected no error for same user, got %v", err)
	}

	// Test permission denied for different user without admin permission
	otherUserID := primitive.NewObjectID()
	otherUserContext := auth.PermissionContext{
		Permissions: []auth.Permission{},
		UserID:      &otherUserID,
	}

	_, err = statsService.GetUserStats(userID, otherUserContext)
	if err != constants.ErrorPermissionDenied {
		t.Errorf("Expected ErrorPermissionDenied, got %v", err)
	}

	// Test repository error
	mockRepo.SetError(true, "repository error")
	_, err = statsService.GetUserStats(userID, adminContext)
	if err != constants.ErrorInternalServerError {
		t.Errorf("Expected ErrorInternalServerError, got %v", err)
	}
}

func testUpdateUserStats(t *testing.T) {
	statsService, mockRepo := createTestUserStatsService()

	userID := primitive.NewObjectID()
	username := "testuser"

	// Create test stats
	stats := user.NewStatsDocument(userID, username)
	mockRepo.CreateStats(stats)

	// Update some fields
	stats.SetField("tasksCreated", 5)
	stats.SetField("tasksCompleted", 3)

	adminContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	// Test successful update
	err := statsService.UpdateUserStats(stats, adminContext)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test permission denied for different user
	otherUserID := primitive.NewObjectID()
	otherUserContext := auth.PermissionContext{
		Permissions: []auth.Permission{},
		UserID:      &otherUserID,
	}

	err = statsService.UpdateUserStats(stats, otherUserContext)
	if err != constants.ErrorPermissionDenied {
		t.Errorf("Expected ErrorPermissionDenied, got %v", err)
	}

	// Test repository error
	mockRepo.SetError(true, "repository error")
	err = statsService.UpdateUserStats(stats, adminContext)
	if err != constants.ErrorInternalServerError {
		t.Errorf("Expected ErrorInternalServerError, got %v", err)
	}
}

func testRecordLogin(t *testing.T) {
	statsService, mockRepo := createTestUserStatsService()

	userID := primitive.NewObjectID()
	username := "testuser"

	// Create test stats
	stats := user.NewStatsDocument(userID, username)
	mockRepo.CreateStats(stats)

	adminContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	// Get initial login count
	initialStats, _ := statsService.GetUserStats(userID, adminContext)
	initialLogins := initialStats.GetInt("totalLogins")

	// Test successful login recording
	err := statsService.RecordLogin(userID, adminContext)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify login count was incremented
	updatedStats, _ := statsService.GetUserStats(userID, adminContext)
	if updatedStats.GetInt("totalLogins") != initialLogins+1 {
		t.Errorf("Expected login count %d, got %d", initialLogins+1, updatedStats.GetInt("totalLogins"))
	}
}

func testIncrementField(t *testing.T) {
	statsService, mockRepo := createTestUserStatsService()

	userID := primitive.NewObjectID()
	username := "testuser"

	// Create test stats
	stats := user.NewStatsDocument(userID, username)
	mockRepo.CreateStats(stats)

	adminContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	// Test incrementing tasks created
	err := statsService.IncrementField(userID, "tasksCreated", 3, adminContext)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the field was incremented
	updatedStats, _ := statsService.GetUserStats(userID, adminContext)
	if updatedStats.GetInt("tasksCreated") != 3 {
		t.Errorf("Expected tasksCreated to be 3, got %d", updatedStats.GetInt("tasksCreated"))
	}

	// Test incrementing again
	err = statsService.IncrementField(userID, "tasksCreated", 2, adminContext)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	updatedStats, _ = statsService.GetUserStats(userID, adminContext)
	if updatedStats.GetInt("tasksCreated") != 5 {
		t.Errorf("Expected tasksCreated to be 5, got %d", updatedStats.GetInt("tasksCreated"))
	}
}

func testIncrementInt64Field(t *testing.T) {
	statsService, mockRepo := createTestUserStatsService()

	userID := primitive.NewObjectID()
	username := "testuser"

	// Create test stats
	stats := user.NewStatsDocument(userID, username)
	mockRepo.CreateStats(stats)

	adminContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	// Test incrementing total time tracked
	err := statsService.IncrementInt64Field(userID, "totalTimeTracked", 3600, adminContext)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the field was incremented
	updatedStats, _ := statsService.GetUserStats(userID, adminContext)
	if updatedStats.GetInt64("totalTimeTracked") != 3600 {
		t.Errorf("Expected totalTimeTracked to be 3600, got %d", updatedStats.GetInt64("totalTimeTracked"))
	}
}

func testUpdateField(t *testing.T) {
	statsService, mockRepo := createTestUserStatsService()

	userID := primitive.NewObjectID()
	username := "testuser"

	// Create test stats
	stats := user.NewStatsDocument(userID, username)
	mockRepo.CreateStats(stats)

	adminContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	// Test updating a field
	err := statsService.UpdateField(userID, "totalLogins", 10, adminContext)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the field was updated
	updatedStats, _ := statsService.GetUserStats(userID, adminContext)
	if updatedStats.GetInt("totalLogins") != 10 {
		t.Errorf("Expected totalLogins to be 10, got %d", updatedStats.GetInt("totalLogins"))
	}
}

func testUpdateFields(t *testing.T) {
	statsService, mockRepo := createTestUserStatsService()

	userID := primitive.NewObjectID()
	username := "testuser"

	// Create test stats
	stats := user.NewStatsDocument(userID, username)
	mockRepo.CreateStats(stats)

	adminContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	// Test updating multiple fields
	fields := map[string]interface{}{
		"totalLogins":  15,
		"tasksCreated": 8,
	}

	err := statsService.UpdateFields(userID, fields, adminContext)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the fields were updated
	updatedStats, _ := statsService.GetUserStats(userID, adminContext)
	if updatedStats.GetInt("totalLogins") != 15 {
		t.Errorf("Expected totalLogins to be 15, got %d", updatedStats.GetInt("totalLogins"))
	}
	if updatedStats.GetInt("tasksCreated") != 8 {
		t.Errorf("Expected tasksCreated to be 8, got %d", updatedStats.GetInt("tasksCreated"))
	}
}

func testPermissionChecks(t *testing.T) {
	statsService, mockRepo := createTestUserStatsService()

	userID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()
	username := "testuser"

	// Create test stats
	stats := user.NewStatsDocument(userID, username)
	mockRepo.CreateStats(stats)

	// Test contexts
	adminContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}

	userContext := auth.PermissionContext{
		Permissions: []auth.Permission{},
		UserID:      &userID,
	}

	otherUserContext := auth.PermissionContext{
		Permissions: []auth.Permission{},
		UserID:      &otherUserID,
	}

	// Test that admin can access any user's stats
	_, err := statsService.GetUserStats(userID, adminContext)
	if err != nil {
		t.Errorf("Admin should be able to access any user's stats, got error: %v", err)
	}

	// Test that user can access their own stats
	_, err = statsService.GetUserStats(userID, userContext)
	if err != nil {
		t.Errorf("User should be able to access their own stats, got error: %v", err)
	}

	// Test that user cannot access other user's stats
	_, err = statsService.GetUserStats(userID, otherUserContext)
	if err != constants.ErrorPermissionDenied {
		t.Errorf("User should not be able to access other user's stats, expected ErrorPermissionDenied, got: %v", err)
	}

	// Test permission checks for other methods
	err = statsService.IncrementField(userID, "tasksCreated", 1, otherUserContext)
	if err != constants.ErrorPermissionDenied {
		t.Errorf("User should not be able to modify other user's stats, expected ErrorPermissionDenied, got: %v", err)
	}

	err = statsService.UpdateField(userID, "totalLogins", 5, otherUserContext)
	if err != constants.ErrorPermissionDenied {
		t.Errorf("User should not be able to modify other user's stats, expected ErrorPermissionDenied, got: %v", err)
	}
}
