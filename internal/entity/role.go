package entity

// Role represent permission roles.
type Role string

//TODO add role description.
const (
	// Admin represent an admin permission role.
	Admin Role = "admin"
	// Operator represent an operator permission role.
	Operator Role = "operator"
	// Viewer represent an viewer permission role.
	Viewer Role = "viewer"
)
