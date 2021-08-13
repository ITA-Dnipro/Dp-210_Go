package entity

type PasswordRestoreReq struct {
	Email string `json:"email"`
}

type PasswordCode struct {
	PasswordRestoreReq
	Code string `json:"code"`
}
