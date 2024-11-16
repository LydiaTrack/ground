package utils

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/LydiaTrack/ground/pkg/domain/email"
	"github.com/LydiaTrack/ground/pkg/registry"
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
	templateName := strings.ToLower(string(emailType))

	// Retrieve the template from the registry
	tmpl, err := registry.GetTemplate(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve template: %w", err)
	}

	// Render the template with the provided data
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return body.String(), nil
}
