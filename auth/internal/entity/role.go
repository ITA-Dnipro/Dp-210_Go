package entity

// Role represent permission roles.
type Role string

const (
	// Admin represent an admin permission role.
	Admin Role = "admin"
	// Operator represent an operator permission role.
	Operator Role = "operator"
	// Doctor represent an doctor permission role.
	Doctor Role = "doctor"
	// Patient represent an patient permission role.
	Patient Role = "patient"
	// Viewer represent an viewer permission role.
	Viewer Role = "viewer"
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
