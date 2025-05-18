package utils

import (
	"errors"
	"regexp"
)

// ValidateUserAvatar checks if the base64 string or is a secure URL
func ValidateUserAvatar(avatar string) error {
	// Check if the string starts with "https://"
	if len(avatar) > 8 && avatar[:8] == "https://" {
		return nil
	}
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
