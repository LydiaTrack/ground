package utils

import (
	"errors"
	"regexp"
)

// ValidateBase64Image checks if the base64 string is a valid image format and within size limits.
func ValidateBase64Image(avatar string) error {
	// Basic base64 regex pattern
	matched, err := regexp.MatchString(`^data:image\/(png|jpg|jpeg);base64,`, avatar)
	if err != nil || !matched {
		return errors.New("invalid avatar format")
	}

	// Check size limit (for example, 1MB)
	if len(avatar) > 1*1024*1024 {
		return errors.New("avatar size exceeds limit")
	}

	return nil
}
