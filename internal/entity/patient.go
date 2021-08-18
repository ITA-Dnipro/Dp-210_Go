package entity

// Patient struct (Model) and request
type Patient struct {
	ID        string `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"omitempty"`
}
