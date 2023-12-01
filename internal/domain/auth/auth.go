package auth

type Permission struct {
	Domain string `json:"domain"`
	Action string `json:"action"`
}

var AdminPermission = Permission{
	Domain: "*",
	Action: "*",
}
