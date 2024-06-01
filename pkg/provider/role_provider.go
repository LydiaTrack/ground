package provider

// RoleProvider for providing roles in specific situations such as creating a new user
type RoleProvider interface {
	// GetDefaultRoles returns the default roles for a new user
	GetDefaultRoleNames() []string
}
