package provider

type SelfRoleProvider struct {
}

func (p SelfRoleProvider) GetDefaultRoleNames() []string {
	return []string{"Lydia Self Service Role"}
}
