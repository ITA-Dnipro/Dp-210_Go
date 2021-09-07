package entity

type Role string

type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	PermissionRole Role   `json:"roles"`
	PasswordHash   []byte `json:"-"`
}

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
