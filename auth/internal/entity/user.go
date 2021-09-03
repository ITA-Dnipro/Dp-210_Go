package entity

// User struct (Model)
type User struct {
	ID string `json:"id"`
	//Name           string `json:"name"`
	Email          string `json:"email"`
	PermissionRole Role   `json:"roles"`
	PasswordHash   []byte `json:"-"`
}

// UserLogin represent new user in request.
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserNewPassword struct {
	UserID          string `json:"id"`
	OldPassword     string `json:"password_old"`
	Password        string `json:"password_new"`
	PasswordConfirm string `json:"password_new_confirm"`
}
