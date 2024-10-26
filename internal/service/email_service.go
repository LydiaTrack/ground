package service

import (
	"fmt"
	"github.com/LydiaTrack/ground/pkg/domain/email"
	"github.com/LydiaTrack/ground/pkg/utils"
	"net/smtp"
)

type SMTPConfig struct {
	Host string
	Port int
}

type SimpleEmailService struct {
	SMTPConfig SMTPConfig
}

// NewSimpleEmailService creates a new SimpleEmailService with the provided SMTP configuration.
func NewSimpleEmailService(smtpConfig SMTPConfig) *SimpleEmailService {
	return &SimpleEmailService{
		SMTPConfig: smtpConfig,
	}
}

// SendEmail sends an email using the provided SendEmailCommand and email type.
// It generates the email body using the specified template.
func (s *SimpleEmailService) SendEmail(command email.SendEmailCommand, emailType string, templateData email.EmailTemplateData) error {
	if emailType == "" {
		emailType = "RESET_PASSWORD"
	}

	// Load the email type configuration from environment variables.
	emailCredentials, err := utils.LoadCredentials(emailType)
	if err != nil {
		return fmt.Errorf("failed to load email type: %v", err)
	}

	// Generate the email body from the template
	body, err := utils.GenerateEmailBody(emailType, templateData)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %v", err)
	}

	// Set up authentication information.
	auth := smtp.PlainAuth("", emailCredentials.Address, emailCredentials.Password, s.SMTPConfig.Host)

	// Create the email message with the generated HTML content.
	msg := []byte(fmt.Sprintf("Subject: %s\r\n", command.Subject) +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		body)

	// Connect to the SMTP server and send the email.
	addr := fmt.Sprintf("%s:%d", s.SMTPConfig.Host, s.SMTPConfig.Port)
	err = smtp.SendMail(addr, auth, emailCredentials.Address, []string{command.To}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully to", command.To)
	return nil
}
