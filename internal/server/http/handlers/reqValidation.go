package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func IsRequestValid(nu interface{}) bool {
	validate := validator.New()
	err := validate.Struct(nu)
	fmt.Println(err)
	return err == nil
}
