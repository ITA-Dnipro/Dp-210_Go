package entity

type PasswordRestoreReq struct {
	Email string `json:"email"`
}

type PasswordCode struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
