package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/LydiaTrack/ground/internal/templates"
	"github.com/LydiaTrack/ground/pkg/domain/email"
)

// LoadCredentialsWithEmailType loads the email address and password from environment variables based on the email type.
// It currently supports "RESET_PASSWORD" as a recognized type.
func LoadCredentialsWithEmailType(emailType email.SupportedEmailType) (email.EmailCredentials, error) {

	addresKey := fmt.Sprintf("EMAIL_TYPE_%s_ADDRESS", strings.ToUpper(string(emailType)))
	passwordKey := fmt.Sprintf("EMAIL_TYPE_%s_PASSWORD", strings.ToUpper(string(emailType)))
	address := os.Getenv(addresKey)
	password := os.Getenv(passwordKey)

	if address == "" || password == "" {
		return email.EmailCredentials{}, fmt.Errorf("environment variables %s or %s not set", addresKey, passwordKey)
	}

	return email.EmailCredentials{
		Address:  address,
		Password: password,
	}, nil
}

// GenerateEmailBodyFromTemplate parses the template and injects the data to generate the final email body.
func GenerateEmailBodyFromTemplate(emailType email.SupportedEmailType, data interface{}) (string, error) {
	// Use lowercase email type as the template name
	tmpl, err := template.ParseFS(templates.FS, fmt.Sprintf("%s.html", strings.ToLower(string(emailType))))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	if tmpl == nil {
		return "", fmt.Errorf("template %s.html not found", strings.ToLower(string(emailType)))
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return body.String(), nil
}
