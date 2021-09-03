package entity

// Role represent permission roles.
type Role string

const (
	Admin    Role = "admin"
	Operator Role = "operator"
	Doctor   Role = "doctor"
	Patient  Role = "patient"
	Viewer   Role = "viewer"
)

func IsAllowedRole(r Role, allowedRoles []Role) bool {
	for _, ar := range allowedRoles {
		if r == ar {
			return true
		}
	}
	return false
}
