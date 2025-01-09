package shared

func GetSignUpEmailTemplate() (string, string) {
	subject := "Orbyters - Sign Up completed"
	body := "Thank you for joining our team! You can now log into your account.\nOrbyters' team."
	return subject, body
}

func GetForgotPasswordEmailTemplate(link string) (string, string) {
	subject := "Orbyters - Forgot password"
	body := "Click the link to reset your password: " + link
	return subject, body
}
