package emails

import (
	"Orbyters/config"
	"fmt"
	"net/smtp"
)

type EmailService struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost:     config.SmtpHost,
		SMTPPort:     config.SmtpPort,
		SMTPUsername: config.SmtpUser,
		SMTPPassword: config.SmtpPass,
		FromEmail:    config.SmtpMail,
	}
}

func (e *EmailService) SendEmail(subject, body string, to []string) error {
	auth := smtp.PlainAuth("", e.SMTPUsername, e.SMTPPassword, e.SMTPHost)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.FromEmail, formatRecipients(to), subject, body))

	serverAddr := fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort)

	fmt.Print(msg, serverAddr)

	err := smtp.SendMail(serverAddr, auth, e.FromEmail, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

func formatRecipients(recipients []string) string {
	return stringJoin(recipients, ",")
}

func stringJoin(elements []string, delimiter string) string {
	result := ""
	for i, element := range elements {
		if i > 0 {
			result += delimiter
		}
		result += element
	}
	return result
}
