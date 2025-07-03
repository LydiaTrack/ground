package handlers

import (
	"net/http"
	"time"

	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/user"
	"github.com/LydiaTrack/ground/pkg/jwt"
	"github.com/LydiaTrack/ground/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatsHandler struct {
	userStatsService *service.UserStatsService
	authService      auth.Service
}

// NewUserStatsHandler creates a new UserStatsHandler
func NewUserStatsHandler(userStatsService *service.UserStatsService, authService auth.Service) UserStatsHandler {
	return UserStatsHandler{
		userStatsService: userStatsService,
		authService:      authService,
	}
}

// GetSelfStats godoc
// @Summary Get current user's stats
// @Description get the current user's statistics including task streak
// @Tags user-stats
// @Accept */*
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /user-stats/me [get]
func (h UserStatsHandler) GetSelfStats(c *gin.Context) {
	userID, err := jwt.ExtractUserIDFromContext(c)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	userId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authContext := auth.PermissionContext{
		UserID: &userId,
	}

	stats, err := h.userStatsService.GetUserStats(userId, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetSelfTaskStreak godoc
// @Summary Get current user's task streak
// @Description get the current user's task streak count with updated values
// @Tags user-stats
// @Accept */*
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /user-stats/me/task-streak [get]
func (h UserStatsHandler) GetSelfTaskStreak(c *gin.Context) {
	userID, err := jwt.ExtractUserIDFromContext(c)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	userId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authContext := auth.PermissionContext{
		UserID: &userId,
	}

	stats, err := h.userStatsService.GetUserStats(userId, authContext)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	// Get current values and calculate accurate streak status
	streakInfo := h.calculateCurrentStreakInfo(stats)

	c.JSON(http.StatusOK, streakInfo)
}

// calculateCurrentStreakInfo calculates the current accurate task streak information
// Similar to updateTaskStreak but for read-only purposes
func (h UserStatsHandler) calculateCurrentStreakInfo(stats user.StatsDocument) map[string]interface{} {
	taskStreak := stats.GetInt("taskStreak")
	tasksCompleted := stats.GetInt("tasksCompleted")
	lastTaskDate := stats.GetTime("lastTaskCompletedDate")

	// Calculate if the streak is still valid (completed task today or yesterday for consecutive days)
	streakValidToday := false
	if !lastTaskDate.IsZero() {
		today := time.Now()
		// Check if last task was completed today
		if lastTaskDate.Year() == today.Year() && lastTaskDate.YearDay() == today.YearDay() {
			streakValidToday = true
		}
	}

	return map[string]interface{}{
		"taskStreak":            taskStreak,
		"tasksCompleted":        tasksCompleted,
		"lastTaskCompletedDate": lastTaskDate,
		"streakValidToday":      streakValidToday,
	}
}
