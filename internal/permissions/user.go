package permissions

import (
	"github.com/LydiaTrack/lydia-base/pkg/auth"
)

var UserCreatePermission = auth.Permission{
	Domain: "user",
	Action: "CREATE",
}

var UserUpdatePermission = auth.Permission{
	Domain: "user",
	Action: "UPDATE",
}

var UserSelfUpdatePermission = auth.Permission{
	Domain: "user",
	Action: "SELF_UPDATE",
}

var UserDeletePermission = auth.Permission{
	Domain: "user",
	Action: "DELETE",
}

var UserReadPermission = auth.Permission{
	Domain: "user",
	Action: "READ",
}
