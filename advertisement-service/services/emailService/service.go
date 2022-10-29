package emailService

import (
	"multiverse/notifier/config"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService interface {
	SendEmail(target, text string) error
}

func NewMailService(senderEmail string) MailService {
	return &SendGridMailService{SenderEmail: senderEmail}
}

type SendGridMailService struct {
	SenderEmail string
}

func (s *SendGridMailService) SendEmail(target, text string) error {
	from := mail.NewEmail("multiverse", s.SenderEmail)
	subject := "welcome email"
	to := mail.NewEmail("for User", target)
	plainTextContent := text
	htmlContent := "<strong>" + text + "</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(config.Configs.Secrets.SendGridToken)
	_, err := client.Send(message)
	return err
}
