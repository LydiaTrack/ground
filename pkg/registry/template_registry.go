package registry

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// TemplateRegistry manages the registration and retrieval of templates.
type TemplateRegistry struct {
	templates map[string]*template.Template
}

// Singleton instance of TemplateRegistry.
var templateRegistry = &TemplateRegistry{
	templates: make(map[string]*template.Template),
}

// RegisterTemplateFromHTML registers a template using raw HTML content.
func RegisterTemplateFromHTML(name, templateContent string) error {
	tmpl, err := template.New(name).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	templateRegistry.templates[name] = tmpl
	return nil
}

// RegisterTemplateFromFile registers a template from a file.
func RegisterTemplateFromFile(name, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", filePath, err)
	}
	return RegisterTemplateFromHTML(name, string(content))
}

// GetTemplate retrieves a template by name.
func GetTemplate(name string) (*template.Template, error) {
	tmpl, exists := templateRegistry.templates[name]
	if !exists {
		return nil, fmt.Errorf("template %s not found", name)
	}
	return tmpl, nil
}

// RenderTemplate renders a template with the provided data and returns the result as a string.
func RenderTemplate(name string, data interface{}) (string, error) {
	tmpl, err := GetTemplate(name)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", name, err)
	}

	return result.String(), nil
}
