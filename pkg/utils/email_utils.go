package utils

import (
	"bytes"
	"fmt"
	"github.com/LydiaTrack/lydia-base/internal/templates"
	"github.com/LydiaTrack/lydia-base/pkg/domain/email"
	"html/template"
	"os"
	"strings"
)

// LoadCredentials loads the email address and password from environment variables based on the email type.
// It currently supports "RESET_PASSWORD" as a recognized type.
func LoadCredentials(emailType string) (email.EmailCredentials, error) {
	var addressEnv, passwordEnv string

	switch emailType {
	case "RESET_PASSWORD":
		addressEnv = "EMAIL_TYPE_RESET_PASSWORD_ADDRESS"
		passwordEnv = "EMAIL_TYPE_RESET_PASSWORD_PASSWORD"
	default:
		return email.EmailCredentials{}, fmt.Errorf("unsupported email type: %s", emailType)
	}

	address := os.Getenv(addressEnv)
	password := os.Getenv(passwordEnv)

	if address == "" || password == "" {
		return email.EmailCredentials{}, fmt.Errorf("environment variables %s or %s not set", addressEnv, passwordEnv)
	}

	return email.EmailCredentials{
		Address:  address,
		Password: password,
	}, nil
}

// GenerateEmailBody parses the template and injects the data to generate the final email body.
func GenerateEmailBody(emailType string, data email.EmailTemplateData) (string, error) {
	// Use lowercase email type as the template name
	tmpl, err := template.ParseFS(templates.FS, fmt.Sprintf("%s.html", strings.ToLower(emailType)))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	if tmpl == nil {
		return "", fmt.Errorf("template %s.html not found", strings.ToLower(emailType))
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return body.String(), nil
}
