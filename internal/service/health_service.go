package service

import (
	"lydia-track-base/internal/domain/health"
)

// GetApplicationHealth returns the health of the application
// By default, it returns UP
func GetApplicationHealth() health.Health {
	return health.Health{Status: "UP"}
}
