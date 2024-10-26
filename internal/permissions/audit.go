package permissions

import (
	"github.com/LydiaTrack/ground/pkg/auth"
)

var AuditCreatePermission = auth.Permission{
	Domain: "audit",
	Action: "CREATE",
}

var AuditDeletePermission = auth.Permission{
	Domain: "audit",
	Action: "DELETE",
}

var AuditReadPermission = auth.Permission{
	Domain: "audit",
	Action: "READ",
}
