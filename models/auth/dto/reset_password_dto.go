package dto

import (
	"github.com/go-playground/validator/v10"
)

type ResetPasswordDto struct {
	Token              string `json:"token"`
	NewPassword        string `json:"newPassword" validate:"required,min=8"`
	ConfirmNewPassword string `json:"confirmNewPassword" validate:"required,min=8"`
}

func (s *ResetPasswordDto) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
