package manager

import (
	"github.com/LydiaTrack/ground/pkg/provider"
)

var roleProviders []provider.RoleProvider

// GetAllDefaultRoleNames returns all default role names from all registered role providers
func GetAllDefaultRoleNames() []string {
	var roleNames []string
	for _, roleProvider := range roleProviders {
		roleNames = append(roleNames, roleProvider.GetDefaultRoleNames()...)
	}
	return roleNames
}

// RegisterRoleProvider registers a new role provider
func RegisterRoleProvider(roleProvider provider.RoleProvider) {
	roleProviders = append(roleProviders, roleProvider)
}
