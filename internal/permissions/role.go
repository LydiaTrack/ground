package permissions

import (
	"github.com/LydiaTrack/lydia-base/pkg/auth"
)

var RoleCreatePermission = auth.Permission{
	Domain: "role",
	Action: "CREATE",
}

var RoleUpdatePermission = auth.Permission{
	Domain: "role",
	Action: "UPDATE",
}

var RoleDeletePermission = auth.Permission{
	Domain: "role",
	Action: "DELETE",
}

var RoleReadPermission = auth.Permission{
	Domain: "role",
	Action: "READ",
}
