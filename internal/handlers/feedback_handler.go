package handlers

import (
	"net/http"

	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/domain/feedback"
	"github.com/LydiaTrack/ground/pkg/utils"
	"github.com/gin-gonic/gin"
)

// FeedbackHandler defines the HTTP handler for feedback-related operations
type FeedbackHandler struct {
	feedbackService service.FeedbackService
}

// NewFeedbackHandler creates a new FeedbackHandler instance
func NewFeedbackHandler(feedbackService service.FeedbackService) FeedbackHandler {
	return FeedbackHandler{feedbackService: feedbackService}
}

// CreateFeedback godoc
// @Summary Create Feedback
// @Description Submit a new feedback.
// @Tags feedback
// @Accept json
// @Produce json
// @Param feedback body feedback.CreateFeedbackCommand true "Feedback data"
// @Success 200 {object} feedback.Feedback
// @Router /feedback [post]
func (h FeedbackHandler) CreateFeedback(c *gin.Context) {
	var cmd feedback.CreateFeedbackCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newFeedback, err := h.feedbackService.CreateFeedback(cmd)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	c.JSON(http.StatusOK, newFeedback)
}

// GetFeedbackByUser godoc
// @Summary Get Feedback by User
// @Description Retrieve all feedback submitted by a specific user.
// @Tags feedback
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} feedback.Feedback
// @Router /feedback/user/{userID} [get]
func (h FeedbackHandler) GetFeedbackByUser(c *gin.Context) {
	userIDParam := c.Param("userID")

	feedbacks, err := h.feedbackService.GetFeedbacksByUser(userIDParam)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

// UpdateFeedbackStatus godoc
// @Summary Update Feedback Status
// @Description Update the status of a specific feedback.
// @Tags feedback
// @Accept json
// @Produce json
// @Param id path string true "Feedback ID"
// @Param status body feedback.FeedbackStatus true "New Status"
// @Success 200 {object} feedback.Feedback
// @Router /feedback/{id}/status [put]
func (h FeedbackHandler) UpdateFeedbackStatus(c *gin.Context) {
	feedbackIDParam := c.Param("id")

	var status struct {
		Status feedback.FeedbackStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.feedbackService.UpdateFeedbackStatus(feedbackIDParam, status.Status)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	updatedFeedback, err := h.feedbackService.GetFeedback(feedbackIDParam)
	if err != nil {
		utils.EvaluateError(err, c)
		return
	}

	c.JSON(http.StatusOK, updatedFeedback)
}
