package provider

type SelfRoleProvider struct {
}

func (p SelfRoleProvider) GetDefaultRoleNames() []string {
	return []string{"Ground Self Service Role"}
}
