package dto

type ResetPasswordDto struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}
