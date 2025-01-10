package dto

import (
	"github.com/go-playground/validator/v10"
)

type SignUpData struct {
	Name            string `json:"name" validate:"required"`
	Surname         string `json:"surname" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

func (s *SignUpData) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
