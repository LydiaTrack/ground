package test

import (
	"fmt"
	"github.com/LydiaTrack/ground/internal/templates"
	"github.com/LydiaTrack/ground/pkg/registry"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/LydiaTrack/ground/internal/service"
	"github.com/LydiaTrack/ground/pkg/domain/email"
	"github.com/LydiaTrack/ground/pkg/domain/resetPassword"
	"github.com/LydiaTrack/ground/pkg/utils"
)

var (
	emailService            *service.SimpleEmailService
	initializedEmailService = false
)

func initializeEmailService() {
	if !initializedEmailService {
		setEmailServiceEnvVariables()

		resetPwSmtp := os.Getenv("EMAIL_TYPE_RESET_PASSWORD_SMTP")
		resetPwPort, err := strconv.Atoi(os.Getenv("EMAIL_TYPE_RESET_PASSWORD_PORT"))
		if err != nil {
			panic(err)
		}
		registerResetPwEmailTemplate()
		smtpConfig := service.SMTPConfig{
			Host: resetPwSmtp,
			Port: resetPwPort,
		}
		// Initialize the SimpleEmailService
		emailService = service.NewSimpleEmailService(smtpConfig)
		initializedEmailService = true
	}
}

// registerResetPwEmailTemplate registers the reset password email template from the embedded FS into the TemplateRegistry.
func registerResetPwEmailTemplate() {
	// Load the template content from the embedded FS
	templateContent, err := templates.FS.ReadFile("reset_password.html")
	if err != nil {
		log.Fatalf("Failed to read reset password template from embedded FS: %v", err)
	}

	// Register the template content in the TemplateRegistry
	err = registry.RegisterTemplateFromHTML("reset_password", string(templateContent))
	if err != nil {
		log.Fatalf("Failed to register reset password email template: %v", err)
	}
}

func setEmailServiceEnvVariables() {
	// Set environment variables for testing
	err := os.Setenv("EMAIL_TYPE_RESET_PASSWORD_ADDRESS", "no-reply@renoten.com")
	if err != nil {
		return
	}
	err = os.Setenv("EMAIL_TYPE_RESET_PASSWORD_PASSWORD", "HFJ3qpj-bxc.uck5fxv")
	if err != nil {
		return
	}
	err = os.Setenv("EMAIL_TYPE_RESET_PASSWORD_SMTP", "smtpout.secureserver.net")

	if err != nil {
		return
	}
	err = os.Setenv("EMAIL_TYPE_RESET_PASSWORD_PORT", "587")
	if err != nil {
		return
	}
}

func TestEmailService(t *testing.T) {
	initializeEmailService()

	t.Run("SendEmailWithResetPasswordType", testSendEmailWithResetPasswordType)
	t.Run("LoadEmailTypeResetPassword", testLoadEmailTypeResetPassword)
}

func testSendEmailWithResetPasswordType(t *testing.T) {

	fmt.Printf("Sending email with reset password type\n")
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
	fmt.Printf("Sending email with code: %s\n", code)
	err = emailService.SendEmail(command, "RESET_PASSWORD", email.TemplateContext{
		Data: resetPassword.EmailTemplateData{
			Code:     code,
			Username: "testUser",
		},
	})
	fmt.Printf("Email sent\n")
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
	emailTypeConfig, err := utils.LoadCredentialsWithEmailType(email.EmailTypeResetPassword)

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
