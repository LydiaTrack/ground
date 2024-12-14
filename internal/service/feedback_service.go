package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/LydiaTrack/ground/internal/log"

	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/email"
	"github.com/LydiaTrack/ground/pkg/domain/feedback"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbackService struct {
	feedbackRepository FeedbackRepository
	userService        UserService
	emailService       SimpleEmailService
}

func NewFeedbackService(feedbackRepository FeedbackRepository, userService UserService) *FeedbackService {
	feedbackSmtp := os.Getenv("EMAIL_TYPE_FEEDBACK_SMTP")
	feedbackPort, err := strconv.Atoi(os.Getenv("EMAIL_TYPE_FEEDBACK_PORT"))
	if err != nil {
		panic(err)
	}

	return &FeedbackService{
		feedbackRepository: feedbackRepository,
		userService:        userService,
		emailService: *NewSimpleEmailService(
			SMTPConfig{
				Host: feedbackSmtp,
				Port: feedbackPort,
			},
		),
	}
}

// FeedbackRepository defines the interface for feedback operations
type FeedbackRepository interface {
	// SaveFeedback saves a feedback record
	SaveFeedback(f feedback.Model) (feedback.Model, error)

	// GetFeedback retrieves a feedback record by ID
	GetFeedback(id primitive.ObjectID) (feedback.Model, error)

	// ExistsFeedback checks if a feedback record exists by ID
	ExistsFeedback(id primitive.ObjectID) (bool, error)

	// GetFeedbacks retrieves all feedback records
	GetFeedbacks() ([]feedback.Model, error)

	// DeleteFeedback deletes a feedback record by ID
	DeleteFeedback(id primitive.ObjectID) error

	// DeleteOlderThan deletes all feedback records older than a specified date
	DeleteOlderThan(date time.Time) error

	// UpdateFeedbackStatus updates the status of a feedback record by ID
	UpdateFeedbackStatus(id primitive.ObjectID, status feedback.FeedbackStatus) error

	// GetFeedbacksByUser retrieves all feedback records submitted by a specific user
	GetFeedbacksByUser(userID primitive.ObjectID) ([]feedback.Model, error)
}

// CreateFeedback creates a new feedback record
func (s FeedbackService) CreateFeedback(command feedback.CreateFeedbackCommand) (feedback.Model, error) {
	f, err := feedback.NewFeedback(
		feedback.WithUserID(command.UserID),
		feedback.WithType(command.Type),
		feedback.WithMessage(command.Message),
	)
	if err != nil {
		return feedback.Model{}, err
	}

	modelAfterCreate, err := s.feedbackRepository.SaveFeedback(*f)
	if err != nil {
		return feedback.Model{}, err
	}

	// After creating and saving the feedback, program should email asynchrously to the given mail address.
	emailDestination := os.Getenv("FEEDBACK_EMAIL_DESTINATION")
	log.Log("Feedback email destination: ", emailDestination)
	if emailDestination != "" {
		go func() {
			err := s.sendFeedbackEmail(emailDestination, modelAfterCreate)
			if err != nil {
				log.Log("Failed to send feedback email: ", err)
			}
		}()
	} else {
		log.Log("Feedback email destination not set")
	}
	log.Log("Feedback created")
	return modelAfterCreate, nil
}

// GetFeedback retrieves a feedback record by ID
func (s FeedbackService) GetFeedback(id string) (feedback.Model, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return feedback.Model{}, err
	}

	model, err := s.feedbackRepository.GetFeedback(objID)
	if err != nil {
		return feedback.Model{}, err
	}

	return model, nil
}

// ExistsFeedback checks if a feedback record exists by ID
func (s FeedbackService) ExistsFeedback(id string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	exists, err := s.feedbackRepository.ExistsFeedback(objID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetFeedbacks retrieves all feedback records
func (s FeedbackService) GetFeedbacks() ([]feedback.Model, error) {
	feedbacks, err := s.feedbackRepository.GetFeedbacks()
	if err != nil {
		return nil, err
	}

	return feedbacks, nil
}

// DeleteFeedback deletes a feedback record by ID
func (s FeedbackService) DeleteFeedback(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.feedbackRepository.DeleteFeedback(objID)
}

// DeleteOlderThan deletes all feedback records older than a specified date
func (s FeedbackService) DeleteOlderThan(date time.Time) error {
	return s.feedbackRepository.DeleteOlderThan(date)
}

// UpdateFeedbackStatus updates the status of a feedback record by ID
func (s FeedbackService) UpdateFeedbackStatus(id string, status feedback.FeedbackStatus) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.feedbackRepository.UpdateFeedbackStatus(objID, status)
}

// GetFeedbacksByUser retrieves all feedback records submitted by a specific user
func (s FeedbackService) GetFeedbacksByUser(userID string) ([]feedback.Model, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	feedbacks, err := s.feedbackRepository.GetFeedbacksByUser(objID)
	if err != nil {
		return nil, err
	}

	return feedbacks, nil
}

// sendFeedbackEmail sends an email notification when new feedback is submitted
func (s FeedbackService) sendFeedbackEmail(emailDestination string, feedbackModel feedback.Model) error {
	// Get the user who submitted the feedback
	authContext := auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	}
	userModel, err := s.userService.GetUser(feedbackModel.UserID.Hex(), authContext)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	subject := "New Feedback Submitted"
	body := fmt.Sprintf("A new feedback has been submitted by User ID: %s.\n\nType: %s\nMessage: %s\nStatus: %s\n\nPlease check the feedbacks page for more information.",
		feedbackModel.UserID.Hex(), feedbackModel.Type, feedbackModel.Message, feedbackModel.Status)
	replyTo := userModel.ContactInfo.Email
	sendMailCmd := email.SendEmailCommand{
		To:      emailDestination,
		Subject: subject,
		Body:    body,
		ReplyTo: &replyTo,
	}

	err = s.emailService.SendEmail(sendMailCmd, email.EmailTypeFeedback, email.TemplateContext{
		Data: feedback.EmailTemplateData{
			ID:        feedbackModel.ID.Hex(),
			UserID:    feedbackModel.UserID.Hex(),
			Username:  userModel.Username,
			Type:      feedbackModel.Type,
			Message:   feedbackModel.Message,
			Status:    feedbackModel.Status,
			CreatedAt: feedbackModel.CreatedAt.Time(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send feedback email: %w", err)
	}

	return nil
}
