package entity

// User struct (Model)
type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	PermissionRole Role   `json:"roles"`
	PasswordHash   []byte `json:"_"`
}

// NewUser represent new user in request.
type NewUser struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}
