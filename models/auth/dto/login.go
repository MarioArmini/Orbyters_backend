package dto

import (
	"github.com/go-playground/validator/v10"
)

type LoginData struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *LoginData) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
