package registry

// RoleProvider defines the interface for providing roles in specific situations.
type RoleProvider interface {
	// GetDefaultRoleNames returns the default roles for a new user.
	GetDefaultRoleNames() []string
}

// roleRegistry is the singleton instance of RoleRegistry.
var roleRegistry = &RoleRegistry{
	roleProviders: []RoleProvider{},
}

// RoleRegistry manages the registration and retrieval of roles.
type RoleRegistry struct {
	roleProviders []RoleProvider
}

// RegisterRoleProvider registers a new RoleProvider in the global RoleRegistry.
func RegisterRoleProvider(roleProvider RoleProvider) {
	roleRegistry.roleProviders = append(roleRegistry.roleProviders, roleProvider)
}

// GetAllDefaultRoleNames retrieves all default role names from all registered RoleProviders in the global RoleRegistry.
func GetAllDefaultRoleNames() []string {
	var roleNames []string
	for _, roleProvider := range roleRegistry.roleProviders {
		roleNames = append(roleNames, roleProvider.GetDefaultRoleNames()...)
	}
	return roleNames
}
