package dto

import (
	"github.com/go-playground/validator/v10"
)

type ForgotPasswordDto struct {
	Email string `json:"email" validate:"required,email"`
}

func (s *ForgotPasswordDto) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
