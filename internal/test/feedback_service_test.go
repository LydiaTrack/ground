package test

import (
	"github.com/LydiaTrack/ground/internal/templates"
	"github.com/LydiaTrack/ground/pkg/registry"
	"log"
	"os"
	"testing"

	"github.com/LydiaTrack/ground/internal/repository"
	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/domain/feedback"
	"github.com/LydiaTrack/ground/pkg/test_support"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	feedbackService     service.FeedbackService
	initializedFeedback = false
)

func initializeFeedbackService() {
	if !initializedFeedback {
		test_support.TestWithMongo()
		repo := repository.GetFeedbackRepository()
		feedbackService = *service.NewFeedbackService(repo)
		registerFeedbackEmailTemplate()
		initializedFeedback = true
	}
}

func setFeedbackEnvVariables() {
	err := os.Setenv("EMAIL_TYPE_RESET_PASSWORD_PORT", "587")
	if err != nil {
		return
	}
	err = os.Setenv("EMAIL_TYPE_FEEDBACK_ADDRESS", "no-reply@renoten.com")
	if err != nil {
		return
	}
	err = os.Setenv("EMAIL_TYPE_FEEDBACK_PASSWORD", "HFJ3qpj-bxc.uck5fxv")
	if err != nil {
		return
	}
	err = os.Setenv("EMAIL_TYPE_FEEDBACK_SMTP", "smtpout.secureserver.net")
	if err != nil {
		return
	}

	err = os.Setenv("EMAIL_TYPE_FEEDBACK_PORT", "587")
	if err != nil {
		return
	}

	err = os.Setenv("FEEDBACK_EMAIL_ACTIVE", "true")
	if err != nil {
		return
	}
	err = os.Setenv("FEEDBACK_EMAIL_DESTINATION", "support@renoten.com")
	if err != nil {
		return
	}
}

// registerFeedbackEmailTemplate registers the reset password email template from the embedded FS into the TemplateRegistry.
func registerFeedbackEmailTemplate() {
	// Load the template content from the embedded FS
	templateContent, err := templates.FS.ReadFile("feedback.html")
	if err != nil {
		log.Fatalf("Failed to read feedback template from embedded FS: %v", err)
	}

	// Register the template content in the TemplateRegistry
	err = registry.RegisterTemplateFromHTML("feedback", string(templateContent))
	if err != nil {
		log.Fatalf("Failed to register feedback email template: %v", err)
	}
}

func TestFeedbackService(t *testing.T) {
	setFeedbackEnvVariables()
	initializeFeedbackService()

	t.Run("CreateFeedback", testCreateFeedback)
	t.Run("GetFeedbackByUser", testGetFeedbackByUser)
}

func testCreateFeedback(t *testing.T) {
	// Create a new feedback command
	command := feedback.CreateFeedbackCommand{
		UserID:  primitive.NewObjectID(),
		Message: "Hi there! I really enjoy using your to-do app to keep my tasks organized. One feature I think would make it even better is the ability to set recurring tasks. For example, I have weekly and monthly tasks that I need to complete regularly, and it would be super helpful if I could set these tasks to repeat automatically at a set interval (like every week or month) instead of having to manually add them each time. Thanks for considering this idea, and keep up the great work!",
		Type:    feedback.FeatureRequest,
	}

	// Create the feedback
	createdFeedback, err := feedbackService.CreateFeedback(command)
	if err != nil {
		t.Errorf("Error creating feedback: %s", err)
	}

	if createdFeedback.Message != command.Message {
		t.Errorf("Expected feedback message: %s, got: %s", command.Message, createdFeedback.Message)
	}

	// Check if the feedback exists
	exists, err := feedbackService.ExistsFeedback(createdFeedback.ID.Hex())
	if err != nil {
		t.Errorf("Error checking if feedback exists: %s", err)
	}

	if !exists {
		t.Errorf("Expected feedback to exist")
	}

	// Retrieve the feedback and verify
	retrievedFeedback, err := feedbackService.GetFeedback(createdFeedback.ID.Hex())
	if err != nil {
		t.Errorf("Error retrieving feedback: %s", err)
	}

	if retrievedFeedback.ID != createdFeedback.ID {
		t.Errorf("Expected feedback ID: %s, got: %s", createdFeedback.ID, retrievedFeedback.ID)
	}
}

func testGetFeedbackByUser(t *testing.T) {
	userID := primitive.NewObjectID()

	// Create feedback entries for a specific user
	command1 := feedback.CreateFeedbackCommand{
		UserID:  userID,
		Message: "First feedback message.",
		Type:    feedback.BugReport,
	}

	command2 := feedback.CreateFeedbackCommand{
		UserID:  userID,
		Message: "Second feedback message.",
		Type:    feedback.General,
	}

	_, err := feedbackService.CreateFeedback(command1)
	if err != nil {
		t.Errorf("Error creating first feedback: %s", err)
	}

	_, err = feedbackService.CreateFeedback(command2)
	if err != nil {
		t.Errorf("Error creating second feedback: %s", err)
	}

	// Retrieve feedbacks by user
	feedbacks, err := feedbackService.GetFeedbacksByUser(userID.Hex())
	if err != nil {
		t.Errorf("Error retrieving feedbacks by user: %s", err)
	}

	if len(feedbacks) != 2 {
		t.Errorf("Expected 2 feedback entries, got: %d", len(feedbacks))
	}

	if feedbacks[0].UserID != userID || feedbacks[1].UserID != userID {
		t.Errorf("Feedback entries do not match the expected user ID")
	}
}
