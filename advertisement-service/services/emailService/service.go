package emailService

import (
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
	mg := mailgun.NewMailgun("sandboxb55f46d9ae404e72bf9ba93b3c12c604.mailgun.org", apiKey)
    m := mg.NewMessage(
        "Excited User <mailgun@YOUR_DOMAIN_NAME>",
        "Hello",
        "Testing some Mailgun awesomeness!",
        "YOU@YOUR_DOMAIN_NAME",
    )

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    _, _, err := mg.Send(ctx, m)
	return err
}
