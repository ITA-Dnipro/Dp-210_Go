package role

// Role represent permission roles.
type Role string

const (
	// Admin represent an admin permission role.
	Admin Role = "admin"
	// Operator represent an operator permission role.
	Operator Role = "operator"
	// Viewer represent an viewer permission role.
	Viewer Role = "viewer"
	// Doctor represent an viewer permission role.
	//Doctor Role = "doctor"
)

// IsAllowedRole check if role is in allowed roles.
func IsAllowedRole(r Role, allowedRoles []Role) bool {
	for _, ar := range allowedRoles {
		if r == ar {
			return true
		}
	}
	return false
}
