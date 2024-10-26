package test

import (
	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/domain/email"
	"github.com/LydiaTrack/ground/pkg/utils"
	"os"
	"testing"
)

var (
	emailService            *service.SimpleEmailService
	initializedEmailService = false
)

func initializeEmailService() {
	if !initializedEmailService {
		// Set up SMTP configuration
		smtpConfig := service.SMTPConfig{
			Host: "smtp.gmail.com",
			Port: 587,
		}
		// Initialize the SimpleEmailService
		emailService = &service.SimpleEmailService{
			SMTPConfig: smtpConfig,
		}
		initializedEmailService = true
	}
}

func TestEmailService(t *testing.T) {
	initializeEmailService()

	t.Run("SendEmailWithResetPasswordType", testSendEmailWithResetPasswordType)
	t.Run("LoadEmailTypeResetPassword", testLoadEmailTypeResetPassword)
}

func testSendEmailWithResetPasswordType(t *testing.T) {
	// Set environment variables for testing
	err := os.Setenv("EMAIL_TYPE_RESET_PASSWORD_ADDRESS", "mcsnturk@gmail.com")
	if err != nil {
		return
	}
	err = os.Setenv("EMAIL_TYPE_RESET_PASSWORD_PASSWORD", "kork yjwi zpch lafj")
	if err != nil {
		return
	}

	// Create a SendEmailCommand
	command := email.SendEmailCommand{
		To:      "muratcansenturk2000@hotmail.com",
		Subject: "Test Email",
		Body:    "<p>This is a test email.</p>",
	}

	// Call the SendEmail method
	code, err := utils.Generate6DigitCode(false)
	if err != nil {
		t.Errorf("Failed to generate code: %v", err)
	}
	err = emailService.SendEmail(command, "RESET_PASSWORD", email.EmailTemplateData{
		Code:     code,
		Username: "testuser",
	})
	if err != nil {
		t.Errorf("Failed to send email: %v", err)
	}
}

func testLoadEmailTypeResetPassword(t *testing.T) {
	// Set environment variables for testing
	err := os.Setenv("EMAIL_TYPE_RESET_PASSWORD_ADDRESS", "test@automatic.com")
	if err != nil {
		return
	}
	err = os.Setenv("EMAIL_TYPE_RESET_PASSWORD_PASSWORD", "password123")
	if err != nil {
		return
	}

	// Load the email type
	emailTypeConfig, err := utils.LoadCredentials("RESET_PASSWORD")

	if err != nil {
		t.Errorf("Failed to load email type: %v", err)
	}

	// Validate the loaded values
	if emailTypeConfig.Address != "test@automatic.com" {
		t.Errorf("Expected 'test@automatic.com', got '%s'", emailTypeConfig.Address)
	}

	if emailTypeConfig.Password != "password123" {
		t.Errorf("Expected 'password123', got '%s'", emailTypeConfig.Password)
	}
}
