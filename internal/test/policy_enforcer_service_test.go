package test

import (
	"lydia-track-base/internal/auth"
	"lydia-track-base/internal/test_support"
	"testing"
)

// TestNewPolicyEnforcerService Create a new policy enforcer service instance
func TestNewPolicyEnforcerService(t *testing.T) {
	test_support.TestWithMongo()
	auth.InitializePolicyEnforcer()
	auth.GetPolicyEnforcer()
}
