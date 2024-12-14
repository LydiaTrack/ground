package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/LydiaTrack/ground/pkg/domain/email"
	"github.com/LydiaTrack/ground/pkg/utils"
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

// GenerateMessageID generates a unique Message-ID for the email.
func GenerateMessageID(domain string) (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	id := base64.URLEncoding.EncodeToString(b)
	return fmt.Sprintf("<%s@%s>", id, domain), nil
}

// SendEmail sends an email using the provided SendEmailCommand and email type.
// It generates the email body using the specified template.
// SendEmail function using strings.Builder
func (s *SimpleEmailService) SendEmail(command email.SendEmailCommand, emailType email.SupportedEmailType, templateData email.TemplateContext) error {
	if emailType == "" {
		return fmt.Errorf("email type is required")
	}

	// Load the email type configuration from environment variables.
	emailCredentials, err := utils.LoadCredentialsWithEmailType(emailType)
	if err != nil {
		return fmt.Errorf("failed to load email type: %v", err)
	}

	// Generate the email body from the template
	body, err := utils.GenerateEmailBodyFromTemplate(emailType, templateData.Data)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %v", err)
	}

	// Generate a unique Message-ID
	messageID, err := GenerateMessageID(strings.Split(emailCredentials.Address, "@")[1])
	if err != nil {
		return fmt.Errorf("failed to generate Message-ID: %v", err)
	}

	// Format the date
	date := time.Now().Format(time.RFC1123Z)

	// Use strings.Builder to construct the email
	var msgBuilder strings.Builder

	msgBuilder.WriteString(fmt.Sprintf("From: %s\r\n", emailCredentials.Address))
	msgBuilder.WriteString(fmt.Sprintf("To: %s\r\n", command.To))
	msgBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", command.Subject))
	msgBuilder.WriteString(fmt.Sprintf("Message-ID: %s\r\n", messageID))
	msgBuilder.WriteString(fmt.Sprintf("Date: %s\r\n", date))
	msgBuilder.WriteString("MIME-Version: 1.0\r\n")
	msgBuilder.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")

	// Add Reply-To if specified
	if command.ReplyTo != nil {
		msgBuilder.WriteString(fmt.Sprintf("Reply-To: %s\r\n", *command.ReplyTo))
	}

	// Add a blank line to separate headers from the body
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(body)

	// Convert the built string to a byte slice
	msg := []byte(msgBuilder.String())

	// Set up authentication information.
	auth := smtp.PlainAuth("", emailCredentials.Address, emailCredentials.Password, s.SMTPConfig.Host)

	// Connect to the SMTP server and send the email.
	addr := fmt.Sprintf("%s:%d", s.SMTPConfig.Host, s.SMTPConfig.Port)
	err = smtp.SendMail(addr, auth, emailCredentials.Address, []string{command.To}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully to", command.To)
	return nil
}
