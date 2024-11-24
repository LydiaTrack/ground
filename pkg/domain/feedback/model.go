package feedback

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FeedbackType represents the type of feedback (e.g., Feature Request, Bug Report, General)
type FeedbackType string

// FeedbackStatus represents the status of a feedback
type FeedbackStatus string

const (
	// Feedback Types
	FeatureRequest FeedbackType = "featureRequest"
	BugReport      FeedbackType = "bugReport"
	General        FeedbackType = "general"

	// Feedback Statuses
	Submitted      FeedbackStatus = "submitted"
	UnderReview    FeedbackStatus = "underReview"
	InProgress     FeedbackStatus = "inProgress"
	Resolved       FeedbackStatus = "resolved"
	Rejected       FeedbackStatus = "rejected"
	AddedToRoadmap FeedbackStatus = "addedToRoadmap"
)

// Model represents a feedback message submitted by a user
type Model struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"userID" bson:"userId"`
	Type      FeedbackType       `json:"type" bson:"type"`
	Message   string             `json:"message" bson:"message"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	Status    FeedbackStatus     `json:"status" bson:"status"`
}

// EmailTemplateData represents the data required to generate an email template
type EmailTemplateData struct {
	ID        string
	UserID    string
	Type      FeedbackType
	Message   string
	Status    FeedbackStatus
	CreatedAt time.Time
}

// Option defines a function type for applying options to Feedback
type Option func(*Model) error

// NewFeedback creates a new Feedback instance using the options pattern
func NewFeedback(opts ...Option) (*Model, error) {
	f := &Model{
		ID:        primitive.NewObjectID(),
		CreatedAt: primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond)),
		Status:    Submitted, // Default status
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(f); err != nil {
			return nil, err
		}
	}

	// Validate required fields
	if err := f.validate(); err != nil {
		return nil, err
	}

	return f, nil
}

// WithUserID sets the UserID field
func WithUserID(userID primitive.ObjectID) Option {
	return func(f *Model) error {
		if userID == primitive.NilObjectID {
			return errors.New("user ID is required")
		}
		f.UserID = userID
		return nil
	}
}

// WithType sets the Type field
func WithType(feedbackType FeedbackType) Option {
	return func(f *Model) error {
		if !isValidFeedbackType(feedbackType) {
			return fmt.Errorf("invalid feedback type: %s", feedbackType)
		}
		f.Type = feedbackType
		return nil
	}
}

// WithMessage sets the Message field
func WithMessage(message string) Option {
	return func(f *Model) error {
		if message == "" {
			return errors.New("message cannot be empty")
		}
		f.Message = message
		return nil
	}
}

// WithStatus sets the Status field (optional)
func WithStatus(status FeedbackStatus) Option {
	return func(f *Model) error {
		if !isValidFeedbackStatus(status) {
			return fmt.Errorf("invalid feedback status: %s", status)
		}
		f.Status = status
		return nil
	}
}

// validate checks if the required fields are set
func (f *Model) validate() error {
	if f.UserID == primitive.NilObjectID {
		return errors.New("user ID is required")
	}
	if f.Message == "" {
		return errors.New("message is required")
	}
	if !isValidFeedbackType(f.Type) {
		return fmt.Errorf("invalid feedback type: %s", f.Type)
	}
	if !isValidFeedbackStatus(f.Status) {
		return fmt.Errorf("invalid feedback status: %s", f.Status)
	}
	return nil
}

// isValidFeedbackType checks if the feedback type is valid
func isValidFeedbackType(feedbackType FeedbackType) bool {
	switch feedbackType {
	case FeatureRequest, BugReport, General:
		return true
	default:
		return false
	}
}

// isValidFeedbackStatus checks if the feedback status is valid
func isValidFeedbackStatus(status FeedbackStatus) bool {
	switch status {
	case Submitted, UnderReview, InProgress, Resolved, Rejected, AddedToRoadmap:
		return true
	default:
		return false
	}
}

// UpdateStatus updates the status of a feedback instance
func (f *Model) UpdateStatus(newStatus FeedbackStatus) error {
	if !isValidFeedbackStatus(newStatus) {
		return fmt.Errorf("invalid feedback status: %s", newStatus)
	}
	f.Status = newStatus
	return nil
}

// FormatForEmail returns a formatted string of the feedback for sending in an email
func (f *Model) FormatForEmail() string {
	return fmt.Sprintf(
		"Feedback ID: %s\nUser ID: %s\nType: %s\nMessage: %s\nStatus: %s\nSubmitted At: %s",
		f.ID, f.UserID, f.Type, f.Message, f.Status, f.CreatedAt.Time().Format(time.RFC3339),
	)
}
