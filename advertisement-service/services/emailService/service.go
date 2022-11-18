package emailService

import (
	"context"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

var yourDomain string = "sandboxb55f46d9ae404e72bf9ba93b3c12c604.mailgun.org"
var privateAPIKey string = "your-private-key"
var senderEmail string = "armingodarzi1380@gmail.com"

type MailService interface {
	SendEmail(target, text string) error
}

func NewMailService() MailService {
	return &MailgunMailService{}
}

type MailgunMailService struct {
}

func (es *MailgunMailService) SendEmail(target, text string) error {
	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	subject := "Informing you about your addvertisement status"

	message := mg.NewMessage(senderEmail, subject, text, target)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)
	log.Println(resp, id)
	return err
}
